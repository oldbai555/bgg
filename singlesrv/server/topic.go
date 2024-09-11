/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package server

import (
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/server/mq"
)

var (
	MqTopicSyncFile     *mq.TopicSt
	MqTopicCacheAllFile *mq.TopicSt
)

func InitTopic() error {
	var err error
	MqTopicSyncFile, err = mq.NewTopicSt(client.MqTopic_MqTopicSyncFile.String(), MqTopicBySyncFileHandler)
	MqTopicCacheAllFile, err = mq.NewTopicSt(client.MqTopic_MqTopicCacheAllFile.String(), MqTopicByCacheAllFileHandler)
	if err != nil {
		return err
	}
	return nil
}
