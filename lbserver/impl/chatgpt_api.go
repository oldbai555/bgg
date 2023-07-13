package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/client/lbchatgpt"
	"github.com/oldbai555/bgg/lbserver/impl/service"
	"github.com/oldbai555/bgg/pkg/gin_tool"
	"github.com/oldbai555/lbtool/log"
)

func ChatCompletion(c *gin.Context) {
	var req lbchatgpt.ChatCompletionReq
	handler := gin_tool.NewHandler(c)

	err := handler.BindAndValidateReq(&req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}

	if err = req.ValidateAll(); err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}

	rsp, err := service.ChatgptServer.ChatCompletion(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
