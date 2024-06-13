/**
 * @Author: zjj
 * @Date: 2024/6/3
 * @Desc:
**/

package mysql

import (
	"context"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
)

func InitDefaultAccount() {
	if User == nil {
		return
	}
	err := User.FirstOrCreate(context.Background(), map[string]interface{}{
		client.FieldUsername_: "oldbai",
		client.FieldPassword_: utils.StrMd5("oldbai"),
		client.FieldNickname_: "大白哥哥",
		client.FieldRole_:     int32(client.ModelUser_RoleAdmin),
	}, &client.ModelUser{})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
