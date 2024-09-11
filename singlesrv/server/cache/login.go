/**
 * @Author: zjj
 * @Date: 2024/3/25
 * @Desc:
**/

package cache

import (
	"github.com/oldbai555/bgg/singlesrv/client"
	"time"
)

func SetLoginInfo(sid string, user *client.BaseUser) error {
	rdb, err := Rdb()
	if err != nil {
		return err
	}

	err = rdb.SetPb(sid, user, time.Hour*24)
	if err != nil {
		return err
	}

	return nil
}

func GetLoginInfo(sid string) (*client.BaseUser, error) {
	var info client.BaseUser
	rdb, err := Rdb()
	if err != nil {
		return nil, err
	}

	err = rdb.GetPb(sid, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func DelLoginInfo(sid string) error {
	rdb, err := Rdb()
	if err != nil {
		return err
	}

	err = rdb.Del(sid)
	if err != nil {
		return err
	}

	return nil
}
