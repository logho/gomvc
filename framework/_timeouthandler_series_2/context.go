package framework_1

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ResponseData struct {
	Status int
	Data   interface{}
}

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	handler        ControllerHandler
	rwMutex        *sync.Mutex
	timeout        bool
}

func NewContext(r *http.Request, rw http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: rw,
		ctx:            r.Context(),
		rwMutex:        &sync.Mutex{},
	}

}

//###  getter & setter for base attr

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *Context) SetHasTimeout() {
	ctx.timeout = true
}

func (ctx *Context) GetRwMutex() *sync.Mutex {
	return ctx.rwMutex
}

func (ctx *Context) isTimeout() bool {
	return ctx.timeout
}

// getter & setter for base attr ###

//### context func

func (ctx *Context) BaseContext() context.Context {
	return ctx.ctx
}

func (ctx *Context) DeadLine() (time.Time, bool) {
	return ctx.ctx.Deadline()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.ctx.Value(key)
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.ctx.Done()
}

func (ctx *Context) Err() error {
	return ctx.ctx.Err()
}

//### get request parse func

func (ctx *Context) IntParam(key string, defaultVal int) int {
	params := ctx.AllParams()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			intparam, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return defaultVal
			}
			return intparam
		}
	}
	return defaultVal
}

func (ctx *Context) StringParam(key string, defaultVal string) string {
	params := ctx.AllParams()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return defaultVal
}

func (ctx *Context) ArrayParam(key string, defaultVal []string) []string {
	params := ctx.AllParams()
	if vals, ok := params[key]; ok {
		return vals
	}
	return defaultVal
}

func (ctx *Context) AllParams() map[string][]string {
	if ctx.request != nil {
		return ctx.request.URL.Query()
	}
	return map[string][]string{}
}

// get request parse func  ###

//### post request parse fun

func (ctx *Context) IntFormParam(key string, defaultVal int) int {
	params := ctx.FormAllParams()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			intparam, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return defaultVal
			}
			return intparam
		}
	}
	return defaultVal
}

func (ctx *Context) StringFormParam(key string, defaultVal string) string {

	params := ctx.FormAllParams()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return defaultVal
}

func (ctx *Context) ArrayFormParam(key string, defaultVal []string) []string {
	params := ctx.FormAllParams()
	if vals, ok := params[key]; ok {
		return vals
	}
	return defaultVal
}

func (ctx *Context) FormAllParams() map[string][]string {
	if ctx.request != nil {
		if err := ctx.request.ParseForm(); err != nil {
			return map[string][]string{}
		}
		return ctx.request.PostForm
	}
	return map[string][]string{}
}

func (ctx *Context) ParseJsonBody(obj interface{}) error {
	if ctx.request != nil {

		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	}
	return errors.New("request must be non nil")

}

// post request parse fun ###

//### response

func (ctx *Context) JsonResp(status int, obj interface{}) error {
	if ctx.isTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	respJson, err := json.Marshal(obj)
	if err != nil {
		ctx.responseWriter.WriteHeader(http.StatusInternalServerError)
		return err
	}
	ctx.responseWriter.WriteHeader(status)
	ctx.responseWriter.Write([]byte(respJson))
	return nil

}

func (ctx *Context) Html(status int, htmlString template.HTML) error {
	if ctx.isTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "text/html;charset=UTF-8")
	ctx.responseWriter.WriteHeader(status)
	ctx.responseWriter.Write([]byte(htmlString))
	return nil

}

func (ctx *Context) text(status int, text string) error {
	if ctx.isTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	ctx.responseWriter.WriteHeader(status)
	ctx.responseWriter.Write([]byte(text))
	return nil
}

// response  ###
