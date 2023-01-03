package main

import (
	"context"
	"fmt"
	"gomvc/framework_1"
	"log"
	"net/http"
	"time"
)

// there is a problem in timeout logic

func HelloControllerHandler(ctx *framework_1.Context) error {
	finishChan := make(chan struct{}, 1)
	panicChan := make(chan struct{}, 1)

	timeoutCtx, cancel := context.WithTimeout(ctx.BaseContext(), 1*time.Second)
	defer cancel()
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- struct{}{}
			}
		}()
		// the service block
		time.Sleep(1 * time.Second)

		ctx.JsonResp(http.StatusOK, "ok")
		finishChan <- struct{}{}
	}()

	select {
	case p := <-panicChan:
		ctx.GetRwMutex().Lock()
		defer ctx.GetRwMutex().Unlock()
		log.Println(p)
		ctx.JsonResp(http.StatusInternalServerError, "some panic occurrence")
	case <-finishChan:
		fmt.Println("normal finished")
	case <-timeoutCtx.Done():
		ctx.GetRwMutex().Lock()
		defer ctx.GetRwMutex().Unlock()
		ctx.JsonResp(http.StatusInternalServerError, "time out")
		ctx.SetHasTimeout()
	}
	return nil
}

func HelloControllerWithTimeoutHandler(ctx *framework_1.Context) error {
	finishChan := make(chan framework_1.ResponseData, 1)
	panicChan := make(chan struct{}, 1)

	timeoutCtx, cancel := context.WithTimeout(ctx.BaseContext(), 1*time.Second)
	defer cancel()
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- struct{}{}
			}
		}()
		// the service block
		time.Sleep(1000 * time.Millisecond)

		finishChan <- framework_1.ResponseData{
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
