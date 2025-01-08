/**
 * @Author: zjj
 * @Date: 2024/6/7
 * @Desc:
**/

package lbsingleserver

import (
	"github.com/oldbai555/bgg/service/lbbase"
	"github.com/oldbai555/bgg/service/lbsingleserver/cache"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/micro/uctx"
)

func CheckAuth(ctx uctx.IUCtx) (*lbbase.BaseUser, error) {
	info, err := cache.GetLoginInfo(ctx.Sid())
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	return info, nil
}
