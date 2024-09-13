/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package mq

import (
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/pkg/marshal"
	"github.com/oldbai555/bgg/singlesrv/client"
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

	val, err := marshal.PbMarshal(obj)
	if err != nil {
		return err
	}

	b, err := encodeMsg(ctx.TraceId(), val)
	if err != nil {
		return err
	}

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

	val, err := marshal.PbMarshal(obj)
	if err != nil {
		return err
	}

	b, err := encodeMsg(ctx.TraceId(), val)
	if err != nil {
		return err
	}

	return producer.DeferredPublish(t.topic, delay, b)
}
