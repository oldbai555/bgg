/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package mq

import (
	"bytes"
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/jsonpb"
)

func EncodeMsg(reqId string, corpId uint32, msg []byte) ([]byte, error) {
	info := &client.NsqMsg{
		ReqId:  reqId,
		CorpId: corpId,
		Data:   msg,
	}
	var buf bytes.Buffer
	err := jsonpb.Marshal(&buf, info)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeMsg(msg []byte) (*client.NsqMsg, error) {
	m := new(client.NsqMsg)
	var reader = bytes.NewReader(msg)
	err := jsonpb.Unmarshal(reader, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func Process(msg *nsq.Message, doLogic func(buf []byte) error) error {
	info, err := DecodeMsg(msg.Body)
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

	err = doLogic(info.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
