package framework_1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Core this core package router  implement the dynamic route
// Core  struct implement by trie tree, it's support dynamic route
type Core struct {
	body            string
	router          map[string]*Tree
	defaultDuration int
}

func NewCore(duration int) *Core {
	router := map[string]*Tree{}

	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{
		router:          router,
		defaultDuration: duration,
	}
}

func (c *Core) Get(url string, handler ControllerHandler, timeout ...int) {
	limit := checkTimeout(timeout...)
	upperUrl := strings.ToUpper(url)
	getRouter := c.router["GET"]
	if err := getRouter.AddRouter(upperUrl, handler, limit); err != nil {
		log.Fatal("add router error: ", err)
	}

}

func (c *Core) Post(url string, handler ControllerHandler, timeout ...int) {
	limit := checkTimeout(timeout...)
	upperUrl := strings.ToUpper(url)
	postRouter := c.router["POST"]
	if err := postRouter.AddRouter(upperUrl, handler, limit); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Put(url string, handler ControllerHandler, timeout ...int) {
	limit := checkTimeout(timeout...)
	upperUrl := strings.ToUpper(url)
	putRouter := c.router["PUT"]
	if err := putRouter.AddRouter(upperUrl, handler, limit); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Delete(url string, handler ControllerHandler, timeout ...int) {
	limit := checkTimeout(timeout...)
	upperUrl := strings.ToUpper(url)
	deleteRouter := c.router["DELETE"]
	if err := deleteRouter.AddRouter(upperUrl, handler, limit); err != nil {
		log.Fatal("add router error: ", err)
	}

}
func checkTimeout(timeout ...int) int {
	if len(timeout) > 1 {
		log.Fatal("num of timeout argument must be 1 or nil")
	}
	var limit int
	if len(timeout) == 1 {
		limit = timeout[0]
	}
	return limit
}

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

func (c *Core) GetHandler(request *http.Request) (ControllerHandler, int) {
	url := request.URL.Path
	method := request.Method
	upperUrl := strings.ToUpper(url)
	upperMethod := strings.ToUpper(method)
	if urlRouter, ok := c.router[upperMethod]; ok {
		return urlRouter.GetHandler(upperUrl)
	}
	return nil, 0
}
func (c *Core) errorBody() string {
	if c.body != "" {
		return c.body
	}
	return "<html><head><title>Timeout</title></head><body><h1>Timeout</h1></body></html>"
}

type timeoutWriter struct {
	w    http.ResponseWriter
	h    http.Header
	wbuf bytes.Buffer
	req  *http.Request

	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
	code        int
}

var ErrHandlerTimeout = errors.New("http: Handler timeout")

func (tw *timeoutWriter) Header() http.Header { return tw.h }

func (tw *timeoutWriter) Write(p []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, ErrHandlerTimeout
	}
	if !tw.wroteHeader {
		tw.writeHeaderLocked(http.StatusOK)
	}
	return tw.wbuf.Write(p)
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}

func (tw *timeoutWriter) writeHeaderLocked(code int) {
	checkWriteHeaderCode(code)

	switch {
	case tw.timedOut:
		return
	case tw.wroteHeader:
		if tw.req != nil {
			log.Println("http header has wrote")
		}
	default:
		tw.wroteHeader = true
		tw.code = code
	}
}

func (tw *timeoutWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.writeHeaderLocked(code)
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	tw := &timeoutWriter{
		w:   response,
		h:   make(http.Header),
		req: request,
	}
	ctx := NewContext(request, response)
	handler, duration := c.GetHandler(request)
	done := make(chan struct{}, 1)
	panicChan := make(chan struct{}, 1)
	if duration == 0 {
		duration = c.defaultDuration
	}
	timeoutCtx, cancel := context.WithTimeout(ctx.BaseContext(), time.Duration(duration*1000000))
	ctx.responseWriter = tw
	defer cancel()
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- struct{}{}
			}
		}()
		handler(ctx)
		close(done)
	}()

	//select {
	//case p := <-panicChan:
	//	ctx.GetRwMutex().Lock()
	//	defer ctx.GetRwMutex().Unlock()
	//	log.Println(p)
	//	ctx.JsonResp(http.StatusInternalServerError, "some panic occurrence")
	//case ret := <-finishChan:
	//	ctx.JsonResp(ret.Status, ret.Data)
	//case <-timeoutCtx.Done():
	//	ctx.GetRwMutex().Lock()
	//	defer ctx.GetRwMutex().Unlock()
	//	ctx.JsonResp(http.StatusInternalServerError, "time out")
	//	ctx.SetHasTimeout()
	//}
	select {
	case p := <-panicChan:
		panic(p)
	case <-done:
		tw.mu.Lock()
		defer tw.mu.Unlock()
		dst := response.Header()
		for k, vv := range tw.h {
			dst[k] = vv
		}
		if !tw.wroteHeader {
			tw.code = http.StatusOK
		}
		response.WriteHeader(tw.code)
		response.Write(tw.wbuf.Bytes())
	case <-timeoutCtx.Done():
		tw.mu.Lock()
		defer tw.mu.Unlock()
		response.WriteHeader(http.StatusServiceUnavailable)
		io.WriteString(response, c.errorBody())
		tw.timedOut = true
	}

}
