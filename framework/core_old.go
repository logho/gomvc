package framework

import (
	"log"
	"net/http"
	"strings"
)

// -------- core old version ------------
// support the way of  method + url get router by the double deep map
// it's not support dynamic route

type OldCore struct {
	router map[string]map[string]ControllerHandler
}

func NewOldCore() *OldCore {
	getRouter := map[string]ControllerHandler{}
	postRouter := map[string]ControllerHandler{}
	putRouter := map[string]ControllerHandler{}
	deleteRouter := map[string]ControllerHandler{}
	router := map[string]map[string]ControllerHandler{}
	router["GET"] = getRouter
	router["POST"] = postRouter
	router["PUT"] = putRouter
	router["DELETE"] = deleteRouter
	return &OldCore{
		router: router,
	}
}

func (c *OldCore) Get(url string, handler ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	getRouter := c.router["GET"]
	getRouter[upperUrl] = handler
}

func (c *OldCore) Post(url string, handler ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	getRouter := c.router["POST"]
	getRouter[upperUrl] = handler
}

func (c *OldCore) Put(url string, handler ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	getRouter := c.router["PUT"]
	getRouter[upperUrl] = handler
}

func (c *OldCore) Delete(url string, handler ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	getRouter := c.router["DELETE"]
	getRouter[upperUrl] = handler
}

func (c *OldCore) Group(prefix string) IGroup {
	//return NewGroup(c, prefix)
	return nil
}

func (c *OldCore) GetHandler(request *http.Request) ControllerHandler {
	url := request.URL.Path
	method := request.Method
	upperUrl := strings.ToUpper(url)
	upperMethod := strings.ToUpper(method)
	if urlRouter, ok := c.router[upperMethod]; ok {
		if handler, ok := urlRouter[upperUrl]; ok {
			return handler
		}
	}
	return nil
}

func (c *OldCore) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request, response)
	handler := c.GetHandler(request)
	if handler == nil {
		return
	}
	log.Println("core router")
	//http.TimeoutHandler(handler,)
	handler(ctx)
}
