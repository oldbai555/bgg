package impl

import "github.com/gin-gonic/gin"

func registerPublicApi(h *gin.Engine) {
	p := h.Group("/public")
	p.GET("/GetFrontUser", GetFrontUser)
	p.POST("/GetCategoryList", GetCategoryList)
	p.POST("/GetArticleList", GetArticleList)
	p.POST("/GetCommentList", GetCommentList)
	p.POST("/AddComment", AddComment)
	p.POST("/GetArticle", GetArticle)
}
