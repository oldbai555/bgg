/**
 * @Author: zjj
 * @Date: 2024/3/25
 * @Desc:
**/

package cache

import (
	"github.com/oldbai555/bgg/service/lbuser"
	"github.com/oldbai555/lbtool/log"
	"time"
)

func SetLoginInfo(k string, user *lbuser.BaseUser) error {
	rdb, err := Rdb()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = rdb.SetPb(k, user, time.Hour*24)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func GetLoginInfo(k string) (*lbuser.BaseUser, error) {
	var info lbuser.BaseUser
	rdb, err := Rdb()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = rdb.GetPb(k, &info)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &info, nil
}

func DelLoginInfo(k string) error {
	rdb, err := Rdb()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = rdb.Del(k)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
