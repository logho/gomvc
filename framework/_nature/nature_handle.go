package _nature

import (
	"encoding/json"
	"gomvc/framework"
	"net/http"
	"strconv"
)

func DoubleStatusHttp(ctx *framework.Context) error {
	ctx.JsonResp(200, "200")
	ctx.JsonResp(500, "500")
	return nil
}

func Test(request *http.Request, responseWriter http.ResponseWriter) {
	//query := request.URL.Query()
	resObj := map[string]interface{}{
		"data": nil,
	}
	responseWriter.Header().Set("Content-Type", "application/json")

	hello := request.PostFormValue("hello")
	helloInt, err := strconv.Atoi(hello)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	resObj["data"] = helloInt
	res, err := json.Marshal(resObj)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.Write([]byte(res))
	responseWriter.WriteHeader(http.StatusOK)

}

func LoginHandler(ctx *framework.Context) error {
	ctx.JsonResp(http.StatusOK, "ok login handler")
	return nil

}

func GroupAddPrefixHandler(ctx *framework.Context) error {
	ctx.JsonResp(http.StatusOK, "/:id delete")
	return nil
}

func GroupIdPutPrefixHandler(ctx *framework.Context) error {
	ctx.JsonResp(http.StatusOK, "/:id put")
	return nil
}
func GroupIdGetPrefixHandler(ctx *framework.Context) error {
	ctx.JsonResp(http.StatusOK, "/:id get")
	return nil
}

func GroupIdPostPrefixHandler(ctx *framework.Context) error {
	ctx.JsonResp(http.StatusOK, "/:id list all")
	return nil
}

func GroupDelPrefixHandler(ctx *framework.Context) error {
	ctx.JsonResp(http.StatusOK, "del")
	return nil
}
func GroupGetPrefixHandler(ctx *framework.Context) error {
	ctx.JsonResp(http.StatusOK, "GroupGetPrefixHandler all")
	return nil
}

func GroupInfoDelPrefixHandler(ctx *framework.Context) error {
	ctx.JsonResp(http.StatusOK, "user info del")
	return nil

}
