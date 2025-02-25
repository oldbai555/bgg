/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package lbossserver

import (
	"github.com/oldbai555/bgg/iface"
	"github.com/oldbai555/bgg/pkg/mq"
	"github.com/oldbai555/bgg/service/lbbase"
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
