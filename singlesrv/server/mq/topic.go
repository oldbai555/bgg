/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package mq

import (
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/jsonpb"
	"github.com/oldbai555/micro/uctx"
	"google.golang.org/protobuf/proto"
	"time"
)

var topicMgr = make(map[string]*TopicSt)

type TopicSt struct {
	topic   string
	handler nsq.Handler
}

func NewTopicSt(topic string, handler nsq.HandlerFunc) (*TopicSt, error) {
	_, ok := topicMgr[topic]
	if ok {
		return nil, client.ErrNsqTopicAlready
	}
	topicMgr[topic] = &TopicSt{topic: topic, handler: handler}
	return topicMgr[topic], nil
}

func (t TopicSt) Pub(ctx uctx.IUCtx, obj proto.Message) error {
	if producer == nil {
		return client.ErrNsqProducerConnectFailure
	}

	err := producer.Ping()
	if err != nil {
		return err
	}

	marshal, err := t.Marshal(obj)
	if err != nil {
		return err
	}

	b, err := encodeMsg(ctx.TraceId(), marshal)
	if err != nil {
		return err
	}

	log.Infof("publish mq msg %s data: %s", t.topic, string(marshal))
	return producer.Publish(t.topic, b)
}

func (t TopicSt) DeferredPublish(ctx uctx.IUCtx, delay time.Duration, obj proto.Message) error {
	if producer == nil {
		return client.ErrNsqProducerConnectFailure
	}

	err := producer.Ping()
	if err != nil {
		return err
	}

	marshal, err := t.Marshal(obj)
	if err != nil {
		return err
	}

	b, err := encodeMsg(ctx.TraceId(), marshal)
	if err != nil {
		return err
	}
	log.Infof("deferred publish mq msg %s data: %s", t.topic, string(marshal))
	return producer.DeferredPublish(t.topic, delay, b)
}

func (t TopicSt) Marshal(obj proto.Message) ([]byte, error) {
	buf, err := jsonpb.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (t TopicSt) Unmarshal(buf []byte, obj proto.Message) error {
	err := jsonpb.Unmarshal(buf, obj)
	if err != nil {
		return err
	}
	return nil
}
