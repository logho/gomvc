package main

import (
	"gomvc/framework_1"
	"gomvc/framework_1/_nature"
)

func registerRouter(core *framework_1.Core) {
	//core.Get("/hello", HelloControllerWithTimeoutHandler)
	//core.Get("/double", _nature.DoubleStatusHttp)
	core.Post("/user/login", _nature.LoginHandler, 1)
	userGroup := core.Group("/user")
	{
		//userInfoGroup := userGroup.Group("/info")
		userGroup.Get("/add", _nature.GroupAddPrefixHandler)
		userGroup.Get("/del", _nature.GroupDelPrefixHandler)
		userGroup.Get("/get/all", _nature.GroupGetPrefixHandler)

		//userInfoGroup.Get("/add", _nature.GroupAddPrefixHandler)
		//userInfoGroup.Get("/del", _nature.GroupInfoDelPrefixHandler)
		//userInfoGroup.Get("/get", _nature.GroupGetPrefixHandler)

		userGroup.Delete("/:id", _nature.GroupAddPrefixHandler)
		userGroup.Put("/:id", _nature.GroupIdPutPrefixHandler)
		userGroup.Get("/:id", _nature.GroupIdGetPrefixHandler)
		userGroup.Get("/list/all", _nature.GroupIdPostPrefixHandler)
	}
}
