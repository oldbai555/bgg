/**
 * @Author: zjj
 * @Date: 2024/12/12
 * @Desc:
**/

package lbsingle

import (
	"github.com/oldbai555/bgg/pkg/brpc"
	"github.com/oldbai555/bgg/service/lbbase"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/micro/uctx"
	"net/http"
)

func CheckAuth(ctx uctx.IUCtx) (*lbbase.BaseUser, error) {
	var resp CheckAuthSysRsp
	err := brpc.DoRequest(ctx, ServerName, CheckAuthSysCMDPath, http.MethodPost, &CheckAuthSysReq{
		Sid: ctx.Sid(),
	}, &resp)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	return resp.User, nil
}
