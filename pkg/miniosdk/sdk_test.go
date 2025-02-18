/**
 * @Author: zjj
 * @Date: 2025/2/24
 * @Desc:
**/

package miniosdk

import (
	"testing"
)

const (
	ep = "192.168.226.6:9000"
	ak = "H7ObnyHeQmDiam3kLCpq"
	sk = "JqBsd8zrUW1E9zdvdGqoH1UfPZ1bk3EdYORDbno5"
	bk = "test"
)

func TestNewClient(t *testing.T) {
	cli, err := NewClient(ep, ak, sk)
	if err != nil {
		t.Logf("err:%v\n", err)
		return
	}
	object, err := cli.PreSignedGetObject(bk, "G-GM.xlsx")
	if err != nil {
		t.Logf("err:%v\n", err)
		return
	}
	t.Log(object)
	return
}
