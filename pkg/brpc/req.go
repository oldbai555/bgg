/**
 * @Author: zjj
 * @Date: 2024/12/12
 * @Desc:
**/

package brpc

import (
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/json"
	"github.com/oldbai555/lbtool/pkg/jsonpb"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/pkg/restysdk"
	"github.com/oldbai555/micro/bconst"
	"github.com/oldbai555/micro/uctx"
	"google.golang.org/protobuf/proto"
	"net/url"
)

type Resp struct {
	Data    string `json:"data"`
	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Hint    string `json:"hint"`
}

func DoRequest(_ uctx.IUCtx, _, path, method string, req, out proto.Message) error {
	var body []byte
	val, err := jsonpb.MarshalToString(req)
	body = []byte(val)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var headers = make(map[string]string)
	headers[bconst.ProtocolType] = bconst.PROTO_TYPE_API_JSON

	// 直接走网关
	var target = fmt.Sprintf("%s://%s:%s/gateway", "http", "127.0.0.1", "20000")
	result, err := url.JoinPath(target, path)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	resp, err := restysdk.NewRequest().SetHeaders(headers).SetBody(body).Execute(method, result)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Infof("do http resp is %s", string(resp.Body()))
	var respBody Resp
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if respBody.ErrCode > 0 {
		return lberr.NewErr(respBody.ErrCode, respBody.ErrMsg)
	}

	err = jsonpb.Unmarshal([]byte(respBody.Data), out)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
