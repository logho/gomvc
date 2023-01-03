package middleware

import (
	"context"
	"gomvc/framework"
	"log"
	"net/http"
	"time"
)

//timeout handler middleware version

func Timeout(d time.Duration) framework.ControllerHandler {
	return func(ctx *framework.Context) error {
		finishChan := make(chan struct{}, 1)
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
			//time.Sleep(1000 * time.Millisecond)
			//
			//finishChan <- ResponseData{
			//	Status: http.StatusOK,
			//	Data:   "hello world",
			//}
			ctx.Next()
		}()

		select {
		case p := <-panicChan:
			ctx.GetRwMutex().Lock()
			defer ctx.GetRwMutex().Unlock()
			log.Println(p)
			ctx.JsonResp(http.StatusInternalServerError, "some panic occurrence")
		case ret := <-finishChan:
			ctx.JsonResp(200, ret)
		case <-timeoutCtx.Done():
			ctx.GetRwMutex().Lock()
			defer ctx.GetRwMutex().Unlock()
			ctx.JsonResp(http.StatusInternalServerError, "time out")
			ctx.SetHasTimeout()
		}

		return nil

	}

}
