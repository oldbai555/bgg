/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package mq

import (
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/singlesrv/server/constant"
	"github.com/oldbai555/lbtool/log"
)

func NewProducer(addr string) (*nsq.Producer, error) {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = constant.MqTimeOut

	p, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	//p.SetLogger(nil, 0)
	err = p.Ping()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	log.Infof("init Producer SUCCESS addr:%s", addr)
	return p, nil
}
