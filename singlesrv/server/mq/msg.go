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
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/uctx"
	"google.golang.org/protobuf/proto"
	"reflect"
)

func encodeMsg(reqId string, msg []byte) ([]byte, error) {
	info := &client.NsqMsg{
		ReqId: reqId,
		Data:  msg,
	}
	buf, err := marshal.PbMarshal(info)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func decodeMsg(msg []byte) (*client.NsqMsg, error) {
	m := new(client.NsqMsg)
	err := marshal.PbUnmarshal(msg, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func Process[M proto.Message](msg *nsq.Message, doLogic func(uCtx uctx.IUCtx, msg M) error) error {
	info, err := decodeMsg(msg.Body)
	if err != nil {
		log.Errorf("process err:%v", err)
		return err
	}

	log.SetLogHint(info.ReqId)

	if msg.Attempts > 3 {
		log.Errorf("exceeding maximum limit")
		msg.Finish()
		return nil
	}

	ctx := uctx.NewBaseUCtx()
	ctx.SetTraceId(info.ReqId)

	var obj M
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		obj = reflect.New(reflect.TypeOf(obj).Elem()).Interface().(M)
	}

	err = marshal.PbUnmarshal(info.Data, obj)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	err = doLogic(ctx, obj)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
