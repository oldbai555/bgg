package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/client/lbbill"
	"github.com/oldbai555/bgg/lbserver/impl/service"
	"github.com/oldbai555/lbtool/log"
)

func registerLbbillApi(h *gin.Engine) {
	// 可以利用反射来映射函数进去
	group := h.Group("bill").Use(RegisterJwt())
	group.POST("DelBill", DelBill)
	group.POST("UpdateBill", UpdateBill)
	group.POST("GetBill", GetBill)
	group.POST("DelBillCategory", DelBillCategory)
	group.POST("UpdateBillCategory", UpdateBillCategory)
	group.POST("AddBill", AddBill)
	group.POST("AddBillCategory", AddBillCategory)
	group.POST("GetBillCategory", GetBillCategory)
	group.POST("GetBillCategoryList", GetBillCategoryList)
	group.POST("GetBillList", GetBillList)
}
func DelBill(c *gin.Context) {
	var req lbbill.DelBillReq
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

	rsp, err := service.BillServer.DelBill(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func UpdateBill(c *gin.Context) {
	var req lbbill.UpdateBillReq
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

	rsp, err := service.BillServer.UpdateBill(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func GetBill(c *gin.Context) {
	var req lbbill.GetBillReq
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

	rsp, err := service.BillServer.GetBill(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func DelBillCategory(c *gin.Context) {
	var req lbbill.DelBillCategoryReq
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

	rsp, err := service.BillServer.DelBillCategory(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func UpdateBillCategory(c *gin.Context) {
	var req lbbill.UpdateBillCategoryReq
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

	rsp, err := service.BillServer.UpdateBillCategory(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func AddBill(c *gin.Context) {
	var req lbbill.AddBillReq
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

	rsp, err := service.BillServer.AddBill(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func AddBillCategory(c *gin.Context) {
	var req lbbill.AddBillCategoryReq
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

	rsp, err := service.BillServer.AddBillCategory(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func GetBillCategory(c *gin.Context) {
	var req lbbill.GetBillCategoryReq
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

	rsp, err := service.BillServer.GetBillCategory(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func GetBillCategoryList(c *gin.Context) {
	var req lbbill.GetBillCategoryListReq
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

	rsp, err := service.BillServer.GetBillCategoryList(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
func GetBillList(c *gin.Context) {
	var req lbbill.GetBillListReq
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

	rsp, err := service.BillServer.GetBillList(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
