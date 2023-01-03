package framework_1

func registerRouter(core *Core) {
	//core.Get("/hello", HelloControllerWithTimeoutHandler)
	//core.Get("/double", _nature.DoubleStatusHttp)
	core.Post("/user/login", LoginHandler, 1001)
	userGroup := core.Group("/user")
	{
		//userInfoGroup := userGroup.Group("/info")
		userGroup.Get("/add", GroupAddPrefixHandler, 1)
		userGroup.Get("/del", GroupDelPrefixHandler)
		userGroup.Get("/get/all", GroupGetPrefixHandler)

		//userInfoGroup.Get("/add", _nature.GroupAddPrefixHandler)
		//userInfoGroup.Get("/del", _nature.GroupInfoDelPrefixHandler)
		//userInfoGroup.Get("/get", _nature.GroupGetPrefixHandler)

		userGroup.Delete("/:id", GroupAddPrefixHandler)
		userGroup.Put("/:id", GroupIdPutPrefixHandler)
		userGroup.Get("/:id", GroupIdGetPrefixHandler)
		userGroup.Get("/list/all", GroupIdPostPrefixHandler)
	}
}
