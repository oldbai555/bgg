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
	"github.com/oldbai555/micro/uctx"
	"google.golang.org/protobuf/proto"
)

func encodeMsg(reqId string, msg []byte) ([]byte, error) {
	info := &client.NsqMsg{
		ReqId: reqId,
		Data:  msg,
	}
	buf, err := proto.Marshal(info)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func decodeMsg(msg []byte) (*client.NsqMsg, error) {
	m := new(client.NsqMsg)
	err := proto.Unmarshal(msg, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func Process(msg *nsq.Message, doLogic func(uCtx uctx.IUCtx, buf []byte) error) error {
	info, err := decodeMsg(msg.Body)
	if err != nil {
		log.Errorf("process err:%v, msg %v", err, msg)
		return err
	}

	log.SetLogHint(info.ReqId)

	if msg.Attempts > 3 {
		log.Errorf("exceeding maximum limit %s", string(info.Data))
		msg.Finish()
		return nil
	}

	log.Infof("process mq msg %s", string(msg.Body))
	ctx := uctx.NewBaseUCtx()
	ctx.SetTraceId(info.ReqId)
	err = doLogic(ctx, info.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
