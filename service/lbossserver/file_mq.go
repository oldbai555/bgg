/**
 * @Author: zjj
 * @Date: 2025/2/26
 * @Desc:
**/

package lbossserver

import (
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/pkg/mq"
	"github.com/oldbai555/bgg/service/lboss"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/uctx"
)

// MqTopicBySyncFileHandler 消息队列-保存文件
func MqTopicBySyncFileHandler(msg *nsq.Message) error {
	return mq.Process[*lboss.MqSyncFile](msg, func(ctx uctx.IUCtx, data *lboss.MqSyncFile) error {
		for _, file := range data.FileList {
			err := saveFileToOrm(ctx, file)
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}
		return nil
	})
}
