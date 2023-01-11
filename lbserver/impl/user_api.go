package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/lbuser"
	"github.com/oldbai555/lbtool/log"
)

func registerLbuserApi(h *gin.Engine) {
	h.POST("/Login", Login)

	// 可以利用反射来映射函数进去
	group := h.Group("user").Use(RegisterJwt())
	group.GET("/Logout", Logout)
	group.GET("/GetLoginUser", GetLoginUser)
}
func Login(c *gin.Context) {
	var req lbuser.LoginReq
	handler := NewHandler(c)

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

	rsp, err := lbuserServer.Login(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func Logout(c *gin.Context) {
	var req lbuser.LogoutReq
	handler := NewHandler(c)

	rsp, err := lbuserServer.Logout(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func GetLoginUser(c *gin.Context) {
	var req lbuser.GetLoginUserReq
	handler := NewHandler(c)

	rsp, err := lbuserServer.GetLoginUser(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
