/**
 * @Author: zjj
 * @Date: 2024/10/22
 * @Desc:
**/

package wxminiprogram

import (
	"github.com/oldbai555/bgg/pkg/constant"
	"github.com/oldbai555/bgg/service/lbbase"
	"github.com/oldbai555/lbtool/pkg/json"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/pkg/restysdk"
	"net/url"
)

func Code2Session(req *lbbase.JsCodeToSessionReq) (*lbbase.JsCodeToSessionRsp, error) {
	path, _ := url.JoinPath(constant.WxMiniProgramPath, "sns", "jscode2session")
	request := restysdk.NewRequest()
	response, err := request.SetQueryParams(map[string]string{
		"appid":      req.Appid,
		"secret":     req.Secret,
		"js_code":    req.JsCode,
		"grant_type": "authorization_code",
	}).Get(path)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	var resp lbbase.JsCodeToSessionRsp
	err = json.Unmarshal(response.Body(), &resp)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	if resp.Errcode > 0 {
		return nil, lberr.NewErr(resp.Errcode, resp.Errmsg)
	}
	return &resp, nil
}
