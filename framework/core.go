package framework

import (
	"log"
	"net/http"
	"strings"
)

// Core this core package router  implement the dynamic route
// Core struct implement by trie tree, it's support dynamic route
type Core struct {
	router map[string]*Tree
}

func NewCore() *Core {
	router := map[string]*Tree{}

	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{
		router: router,
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

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request, response)
	handler, _ := c.GetHandler(request)
	if handler == nil {
		return
	}
	log.Println("core router")

	handler(ctx)
}
