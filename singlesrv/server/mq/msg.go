/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package mq

import (
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/json"
)

type Msg struct {
	ReqId  string
	CorpId uint32
	Data   []byte
}

func EncodeMsg(reqId string, corpId uint32, msg []byte) ([]byte, error) {
	info := &Msg{
		ReqId:  reqId,
		CorpId: corpId,
		Data:   msg,
	}
	b, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func DecodeMsg(msg []byte) (*Msg, error) {
	m := new(Msg)
	err := json.Unmarshal(msg, m)
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

	log.Infof("process mq msg %v", msg)

	err = doLogic(info.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
