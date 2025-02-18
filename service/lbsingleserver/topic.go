/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package lbsingleserver

import (
	"github.com/oldbai555/bgg/service/lbbase"
	"github.com/oldbai555/bgg/service/lbsingleserver/iface"
	"github.com/oldbai555/bgg/service/lbsingleserver/mq"
	"github.com/oldbai555/lbtool/pkg/lberr"
)

var (
	MqTopicSyncFile iface.ITopic
)

func InitTopic() error {
	var err error
	MqTopicSyncFile, err = mq.NewTopicSt(lbbase.MqTopic_MqTopicSyncFile.String(), MqTopicBySyncFileHandler)
	if err != nil {
		return lberr.Wrap(err)
	}
	return nil
}
