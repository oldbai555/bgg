/**
 * @Author: zjj
 * @Date: 2025/2/24
 * @Desc:
**/

package lbsingleserver

import (
	"github.com/oldbai555/bgg/pkg/miniosdk"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"sync"
)

var minIoSDK *miniosdk.Client
var minIoOnce sync.Once

func InitMinio() error {
	minIOConf := syscfg.NewMinIOConf("")
	var err error
	minIoOnce.Do(func() {
		minIoSDK, err = miniosdk.NewClient(
			minIOConf.Endpoint,
			minIOConf.AccessKey,
			minIOConf.SecretAccessKey,
		)
	})
	return err
}
