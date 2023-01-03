package main

import (
	framework_1 "gomvc/framework/_timeouthandler_series_2"
	"net/http"
	"time"
)

func main() {

	core := framework_1.NewCore(1)
	handler := http.TimeoutHandler(core, time.Duration(1), "")
	registerRouter(core)
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	server.ListenAndServe()
}
