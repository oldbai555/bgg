/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package cache

import (
	"github.com/oldbai555/bgg/singlesrv/server/constant"
	"github.com/oldbai555/lbtool/log"
)

func GetFileBySortUrl(sortUrl string) (string, error) {
	rdb, err := Rdb()
	if err != nil {
		return "", err
	}
	val, err := rdb.HGet(constant.FileCachePrefix, sortUrl)
	if err != nil {
		return "", err
	}
	return val, nil
}

func SetFileBySortUrl(sortUrl string, fileJsonStr string) error {
	rdb, err := Rdb()
	if err != nil {
		return err
	}
	err = rdb.HSet(constant.FileCachePrefix, sortUrl, fileJsonStr)
	if err != nil {
		return err
	}
	return nil
}

func DelFileBySortUrl(sortUrl string) error {
	rdb, err := Rdb()
	if err != nil {
		return err
	}
	_, err = rdb.HDel(constant.FileCachePrefix, sortUrl)
	if err != nil {
		return err
	}
	return nil
}

func DelAllFileCache() error {
	rdb, err := Rdb()
	if err != nil {
		return err
	}
	all, err := rdb.HGetAll(constant.FileCachePrefix)
	if err != nil {
		return err
	}
	for sortUrl := range all {
		err = DelFileBySortUrl(sortUrl)
		if err != nil {
			log.Errorf("del all file cache: %s,err:%v", sortUrl, err)
		}
	}
	return nil
}
