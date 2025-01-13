/**
 * @Author: zjj
 * @Date: 2024/3/25
 * @Desc:
**/

package cache

import (
	"github.com/oldbai555/bgg/service/lbbase"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"time"
)

func SetLoginInfo(sid string, user *lbbase.BaseUser) error {
	rdb, err := Rdb()
	if err != nil {
		return lberr.Wrap(err)
	}

	err = rdb.SetPb(sid, user, time.Hour*24)
	if err != nil {
		return lberr.Wrap(err)
	}

	return nil
}

func GetLoginInfo(sid string) (*lbbase.BaseUser, error) {
	var info lbbase.BaseUser
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
		return lberr.Wrap(err)
	}

	err = rdb.Del(sid)
	if err != nil {
		return lberr.Wrap(err)
	}

	return nil
}
