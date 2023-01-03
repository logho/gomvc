package framework

import (
	"context"
	"log"
	"net/http"
	"time"
)

func TimeoutHandler(handler ControllerHandler, d time.Duration) ControllerHandler {
	return func(ctx *Context) error {
		finishChan := make(chan ResponseData, 1)
		panicChan := make(chan struct{}, 1)

		timeoutCtx, cancel := context.WithTimeout(ctx.BaseContext(), d)
		defer cancel()
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- struct{}{}
				}
			}()
			// the service block
			time.Sleep(1000 * time.Millisecond)

			handler(ctx)

			finishChan <- ResponseData{
				Status: http.StatusOK,
				Data:   "hello world",
			}
		}()

		select {
		case p := <-panicChan:
			ctx.GetRwMutex().Lock()
			defer ctx.GetRwMutex().Unlock()
			log.Println(p)
			ctx.JsonResp(http.StatusInternalServerError, "some panic occurrence")
		case ret := <-finishChan:
			ctx.JsonResp(ret.Status, ret.Data)
		case <-timeoutCtx.Done():
			ctx.GetRwMutex().Lock()
			defer ctx.GetRwMutex().Unlock()
			ctx.JsonResp(http.StatusInternalServerError, "time out")
			ctx.SetHasTimeout()
		}

		return nil

	}

}
