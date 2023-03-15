package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/client/lbblog"
	"github.com/oldbai555/bgg/lbserver/impl/service"
	"github.com/oldbai555/lbtool/log"
)

func registerLbblogApi(h *gin.Engine) {
	group := h.Group("blog")
	group.Use(RegisterJwt())

	// 可以利用反射来映射函数进去
	group.POST("/GetArticle", GetArticle)
	group.POST("/GetArticleList", GetArticleList)
	group.POST("/AddArticle", AddArticle)
	group.POST("/UpdateArticle", UpdateArticle)
	group.POST("/DelArticle", DelArticle)

	group.POST("/GetCategory", GetCategory)
	group.POST("/GetCategoryList", GetCategoryList)
	group.POST("/AddCategory", AddCategory)
	group.POST("/UpdateCategory", UpdateCategory)
	group.POST("/DelCategory", DelCategory)

	group.POST("/GetComment", GetComment)
	group.POST("/GetCommentList", GetCommentList)
	group.POST("/AddComment", AddComment)
	group.POST("/UpdateComment", UpdateComment)
	group.POST("/DelComment", DelComment)
}

func GetArticleList(c *gin.Context) {
	var req lbblog.GetArticleListReq
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

	rsp, err := service.BlogServer.GetArticleList(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func GetArticle(c *gin.Context) {
	var req lbblog.GetArticleReq
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

	rsp, err := service.BlogServer.GetArticle(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func UpdateArticle(c *gin.Context) {
	var req lbblog.UpdateArticleReq
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

	rsp, err := service.BlogServer.UpdateArticle(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func DelArticle(c *gin.Context) {
	var req lbblog.DelArticleReq
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

	rsp, err := service.BlogServer.DelArticle(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func AddArticle(c *gin.Context) {
	var req lbblog.AddArticleReq
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

	rsp, err := service.BlogServer.AddArticle(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func GetCategoryList(c *gin.Context) {
	var req lbblog.GetCategoryListReq
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

	rsp, err := service.BlogServer.GetCategoryList(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func GetCategory(c *gin.Context) {
	var req lbblog.GetCategoryReq
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

	rsp, err := service.BlogServer.GetCategory(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func UpdateCategory(c *gin.Context) {
	var req lbblog.UpdateCategoryReq
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

	rsp, err := service.BlogServer.UpdateCategory(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func DelCategory(c *gin.Context) {
	var req lbblog.DelCategoryReq
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

	rsp, err := service.BlogServer.DelCategory(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func AddCategory(c *gin.Context) {
	var req lbblog.AddCategoryReq
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

	rsp, err := service.BlogServer.AddCategory(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func GetCommentList(c *gin.Context) {
	var req lbblog.GetCommentListReq
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

	rsp, err := service.BlogServer.GetCommentList(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func GetComment(c *gin.Context) {
	var req lbblog.GetCommentReq
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

	rsp, err := service.BlogServer.GetComment(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func UpdateComment(c *gin.Context) {
	var req lbblog.UpdateCommentReq
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

	rsp, err := service.BlogServer.UpdateComment(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func DelComment(c *gin.Context) {
	var req lbblog.DelCommentReq
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

	rsp, err := service.BlogServer.DelComment(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}

func AddComment(c *gin.Context) {
	var req lbblog.AddCommentReq
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

	rsp, err := service.BlogServer.AddComment(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
	handler.Success(rsp)
}
