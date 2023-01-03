package framework_1

import (
	"net/http"
	"testing"
)

func TestTimeout(t *testing.T) {

	core := NewCore(1000)
	//handler := http.TimeoutHandler(core, time.Duration(1), "")
	registerRouter(core)
	server := &http.Server{
		Addr:    ":8080",
		Handler: core,
	}
	server.ListenAndServe()
}
