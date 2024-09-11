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

type Consumer struct {
	// 监听地址
	addressList []string
	// channel 消费队列管道
	channel string
	// entries topic 和 处理方法或函数
	// clis 消费者列表
	clis []*nsq.Consumer
	// cfg 配制
	cfg *nsq.Config
}

func NewConsumer(channel string, addressList ...string) *Consumer {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = constant.MqTimeOut
	o := &Consumer{
		channel:     channel,
		cfg:         cfg,
		addressList: addressList,
	}
	return o
}

func (n *Consumer) Start() (err error) {
	for _, e := range topicMgr {
		var cli *nsq.Consumer
		cli, err = nsq.NewConsumer(e.topic, n.channel, n.cfg)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
		cli.AddHandler(e.handler)
		cli.SetLogger(nil, 0)
		// 直接连接 nsqd 不走服务发现
		err = cli.ConnectToNSQDs(n.addressList)
		if err != nil {
			return
		}
		n.clis = append(n.clis, cli)
		log.Infof("NewConsumer Success Topic:%s, channel:%s, addressList:%v ", e.topic, n.channel, n.addressList)
	}
	return
}

func (n *Consumer) Stop() {
	for _, cli := range n.clis {
		cli.Stop()
	}
	for _, cli := range n.clis {
		<-cli.StopChan
	}
}
