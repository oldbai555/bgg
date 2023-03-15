package impl

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/bgg/client/lbuser"
	"github.com/oldbai555/bgg/lbserver/impl/cache"
	"github.com/oldbai555/bgg/lbserver/impl/constant"
	"github.com/oldbai555/bgg/pkg/webtool"

	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/pkg/result"
	"github.com/oldbai555/lbtool/utils"
	"github.com/pkg/errors"

	"net/http"
)

const (
	HttpHeaderAuthorization     = "Authorization"
	HttpHeaderContentType       = "Content-Type"
	HttpHeaderUserAgent         = "User-Agent"
	HttpHeaderContentTypeByJson = "application/json"
	defaultRspMsg               = "ok"
)

type Handler struct {
	C *gin.Context
}

func NewHandler(c *gin.Context) *Handler {
	handler := &Handler{
		C: c,
	}
	return handler
}

// BindAndValidateReq 绑定并校验请求参数 - 请求体
// req 必须是指针
func (r *Handler) BindAndValidateReq(req interface{}) error {
	err := r.C.ShouldBindJSON(req)
	if err != nil {
		return err
	}
	return nil
}

func (r *Handler) GetHeader(key string) string {
	return r.C.GetHeader(key)
}

func (r *Handler) GetQuery(key string) (string, bool) {
	return r.C.GetQuery(key)
}

func (r *Handler) CheckUser() (*lbuser.ModelUser, error) {
	claims, err := r.GetClaims()
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	var user lbuser.ModelUser
	err = cache.Rdb.GetJson(context.TODO(), fmt.Sprintf("%d", claims.UserId), &user)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *Handler) GetClaims() (*webtool.Claims, error) {
	claims, ok := r.C.Value(constant.CtxWithClaim).(*webtool.Claims)
	if !ok {
		return nil, constant.ErrNotLoginInfo
	}
	return claims, nil
}

// Response 自定义自定义数据
func (r *Handler) Response(httpCode int, errorCode int, data interface{}, msg string) {
	r.C.Header(HttpHeaderContentType, HttpHeaderContentTypeByJson)
	hint := r.C.Value(constant.LogWithHint)
	jsonResult := result.JSONResult{
		Code:    errorCode,
		Message: msg,
		Data:    data,
		Hint:    fmt.Sprintf("%s", hint),
	}
	log.Infof("jsonRes is %v", jsonResult)
	r.C.JSON(httpCode, jsonResult)
}

// Success 响应数据
func (r *Handler) Success(data interface{}) {
	// 特殊转化一下json
	m, err := utils.StructToMapV2(data)
	if err != nil {
		log.Errorf("err is %v", err)
	}
	log.Infof("val is %v", m)
	r.Response(http.StatusOK, int(lbconst.ErrCode_Success), m, defaultRspMsg)
}

// RespError 响应错误
func (r *Handler) RespError(err error) {
	// 获取根错误
	rootErr := errors.Cause(err)

	if e, ok := rootErr.(*lberr.LbErr); ok {
		r.Response(http.StatusOK, int(e.Code()), nil, e.Message())
		return
	}

	r.Response(http.StatusOK, int(lbconst.ErrCode_ErrAnInvalidRsp), nil, err.Error())
}
