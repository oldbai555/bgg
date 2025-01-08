/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package mq

import (
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/pkg/constant"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/lbtool/pkg/lberr"
)

var producer *nsq.Producer
var consumer *Consumer

func Start() error {
	var err error
	producer, err = NewProducer(constant.MqAddress)
	if err != nil {
		return lberr.Wrap(err)
	}

	consumer = NewConsumer(lbsingle.ServerName, constant.MqAddress)
	err = consumer.Start()
	if err != nil {
		return lberr.Wrap(err)
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
