/**
 * @Author: zjj
 * @Date: 2024/6/3
 * @Desc:
**/

package tool

import (
	"github.com/oldbai555/lbtool/log"
	"net"
)

func GetOnePort() uint32 {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		log.Errorf("获取端口失败:%s", err)
		return 0
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Errorf("监听端口失败:%s", err)
		return 0
	}
	err = l.Close()
	if err != nil {
		log.Errorf("结束监听端口失败:%s", err)
		return 0
	}
	onePort := l.Addr().(*net.TCPAddr).Port
	log.Infof("获取端口成功:%v", onePort)
	return uint32(onePort)
}
