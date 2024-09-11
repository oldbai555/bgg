/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package mq

import (
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/server/constant"
	"github.com/oldbai555/lbtool/log"
)

var producer *nsq.Producer
var consumer *Consumer

func Start() error {
	var err error
	producer, err = NewProducer(constant.MqAddress)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	consumer = NewConsumer(client.ServerName, constant.MqAddress)
	err = consumer.Start()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return err
}

func Stop() {
	if producer != nil {
		producer.Stop()
	}
	if consumer != nil {
		consumer.Stop()
	}
}
