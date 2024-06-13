/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package mq

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/jsonpb"
	"github.com/oldbai555/micro/uctx"
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

	b, err := EncodeMsg(ctx.TraceId(), 0, marshal)
	if err != nil {
		return err
	}

	log.Infof("publish mq msg %s data: %s", t.topic, string(marshal))
	return producer.Publish(t.topic, b)
}

func (t TopicSt) Marshal(obj proto.Message) ([]byte, error) {
	var buf bytes.Buffer
	err := jsonpb.Marshal(&buf, obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t TopicSt) Unmarshal(buf []byte, obj proto.Message) error {
	var reader = bytes.NewReader(buf)
	err := jsonpb.Unmarshal(reader, obj)
	if err != nil {
		return err
	}
	return nil
}
