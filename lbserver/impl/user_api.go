package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/client/lbuser"

	"github.com/oldbai555/lbtool/log"
)

func registerLbuserApi(h *gin.Engine) {
	h.POST("/Login", Login)

	// 可以利用反射来映射函数进去
	group := h.Group("user").Use(RegisterJwt())
	group.GET("/Logout", Logout)
	group.GET("/GetLoginUser", GetLoginUser)
	group.POST("/UpdateLoginUserInfo", UpdateLoginUserInfo)
	group.POST("/DelUser", DelUser)
	group.POST("/AddUser", AddUser)
	group.POST("/GetUserList", GetUserList)
	group.POST("/GetUser", GetUser)
	group.POST("/UpdateUserNameWithRole", UpdateUserNameWithRole)
	group.POST("/ResetPassword", ResetPassword)
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
func UpdateLoginUserInfo(c *gin.Context) {
	var req lbuser.UpdateLoginUserInfoReq
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

	rsp, err := lbuserServer.UpdateLoginUserInfo(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func DelUser(c *gin.Context) {
	var req lbuser.DelUserReq
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

	rsp, err := lbuserServer.DelUser(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func AddUser(c *gin.Context) {
	var req lbuser.AddUserReq
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

	rsp, err := lbuserServer.AddUser(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func GetUserList(c *gin.Context) {
	var req lbuser.GetUserListReq
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

	rsp, err := lbuserServer.GetUserList(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func GetUser(c *gin.Context) {
	var req lbuser.GetUserReq
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

	rsp, err := lbuserServer.GetUser(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func UpdateUserNameWithRole(c *gin.Context) {
	var req lbuser.UpdateUserNameWithRoleReq
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

	rsp, err := lbuserServer.UpdateUserNameWithRole(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func ResetPassword(c *gin.Context) {
	var req lbuser.ResetPasswordReq
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

	rsp, err := lbuserServer.ResetPassword(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func GetFrontUser(c *gin.Context) {
	var req lbuser.GetFrontUserReq
	handler := NewHandler(c)

	rsp, err := lbuserServer.GetFrontUser(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
