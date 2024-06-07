/**
 * @Author: zjj
 * @Date: 2024/3/25
 * @Desc:
**/

package cache

import (
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/lbtool/log"
	"time"
)

func SetLoginInfo(sid string, user *client.BaseUser) error {
	rdb, err := Rdb()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = rdb.SetPb(sid, user, time.Hour*24)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func GetLoginInfo(sid string) (*client.BaseUser, error) {
	var info client.BaseUser
	rdb, err := Rdb()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = rdb.GetPb(sid, &info)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &info, nil
}

func DelLoginInfo(sid string) error {
	rdb, err := Rdb()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = rdb.Del(sid)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
