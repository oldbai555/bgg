/**
 * @Author: zjj
 * @Date: 2024/6/7
 * @Desc:
**/

package lbsingleserver

import (
	"github.com/oldbai555/bgg/service/lbbase"
	"github.com/oldbai555/bgg/service/lbsingleserver/cache"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/uctx"
)

func CheckAuth(ctx uctx.IUCtx) (*lbbase.BaseUser, error) {
	info, err := cache.GetLoginInfo(ctx.Sid())
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return info, nil
}
