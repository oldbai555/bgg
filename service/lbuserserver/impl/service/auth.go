/**
 * @Author: zjj
 * @Date: 2024/3/23
 * @Desc:
**/

package service

import (
	"context"
	"github.com/oldbai555/bgg/service/lbuser"
	"github.com/oldbai555/lbtool/log"
)

func (a *LbuserServer) CheckLoginUser(ctx context.Context, sid string) (interface{}, error) {
	userSysRsp, err := a.GetLoginUser(ctx, &lbuser.GetLoginUserReq{
		Sid: sid,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return userSysRsp.BaseUser, nil
}
