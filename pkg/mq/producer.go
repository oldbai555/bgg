/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package mq

import (
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/pkg/constant"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
)

func NewProducer(addr string) (*nsq.Producer, error) {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = constant.MqTimeOut

	p, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	//p.SetLogger(nil, 0)
	err = p.Ping()
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	log.Infof("init Producer SUCCESS addr:%s", addr)
	return p, nil
}
