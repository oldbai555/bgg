/**
 * @Author: zjj
 * @Date: 2024/8/22
 * @Desc:
**/

package iface

type IWsConn interface {
	GetConnId() string
	GetUid() uint64
	IsLogin() bool
}
