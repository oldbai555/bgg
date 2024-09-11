/**
 * @Author: zjj
 * @Date: 2024/6/7
 * @Desc:
**/

package server

import (
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/server/cache"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/uctx"
)

func CheckAuth(ctx uctx.IUCtx) (*client.BaseUser, error) {
	info, err := cache.GetLoginInfo(ctx.Sid())
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return info, nil
}
