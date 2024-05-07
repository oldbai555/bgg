/**
 * @Author: zjj
 * @Date: 2024/5/7
 * @Desc:
**/

package ginhelper

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/micro/bconst"
	"net/http"
)

type TemplateHelper struct {
	c *gin.Context
}

func NewTemplateHelper(c *gin.Context) *TemplateHelper {
	return &TemplateHelper{
		c: c,
	}
}

func (h *TemplateHelper) Error(err error) {
	hint := h.c.Value(bconst.LogWithHint)
	h.c.HTML(http.StatusOK, "err.html", gin.H{
		"hint": hint,
		"err":  err.Error(),
	})
}
