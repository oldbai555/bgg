/**
 * @Author: zjj
 * @Date: 2024/10/24
 * @Desc:
**/

package cache

import (
	"fmt"
	"github.com/oldbai555/bgg/singlesrv/server/constant"
	"github.com/oldbai555/lbtool/log"
	"time"
)

func SetMpSmsCode(k string, v string) error {
	rdb, err := Rdb()
	if err != nil {
		return err
	}
	err = rdb.Set(fmt.Sprintf("%s_%s", constant.MpSmsCachePrefix, k), []byte(v), 60*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func GetMpSmsCode(k string) (string, error) {
	rdb, err := Rdb()
	if err != nil {
		return "", err
	}
	v, err := rdb.Get(fmt.Sprintf("%s_%s", constant.MpSmsCachePrefix, k))
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func DelMpSmsCode(k string) {
	rdb, err := Rdb()
	if err != nil {
		return
	}
	err = rdb.Del(fmt.Sprintf("%s_%s", constant.MpSmsCachePrefix, k))
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	return
}
