/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package server

import (
	"github.com/oldbai555/bgg/singlesrv/server/constant"
	"github.com/oldbai555/bgg/singlesrv/server/mq"
)

var (
	MqTopicBySyncFile *mq.TopicSt
)

func InitTopic() error {
	var err error
	MqTopicBySyncFile, err = mq.NewTopicSt(constant.MqTopicBySyncFile, MqTopicBySyncFileHandler)
	if err != nil {
		return err
	}
	return nil
}
