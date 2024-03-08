package bhttp

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/oldbai555/bgg/internal/_const"
	"github.com/oldbai555/bgg/internal/bgrpc/tool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/dispatch"
	"github.com/oldbai555/lbtool/pkg/json"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/pkg/restysdk"
	"net/url"
)

type Resp struct {
	Data    string `json:"data"`
	Errcode int32  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Hint    string `json:"hint"`
}

func DoRequest(ctx context.Context, srv, path, method string, protocolType string, req, out proto.Message) error {
	d, err := tool.New()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	node, err := dispatch.Route(ctx, d, srv)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var body []byte
	switch protocolType {
	case _const.PROTO_TYPE_API_JSON:
		m := jsonpb.Marshaler{
			EmitDefaults: true,
			OrigName:     true,
		}
		var val string
		val, err = m.MarshalToString(req)
		body = []byte(val)
	case _const.PROTO_TYPE_PROTO3:
		body, err = proto.Marshal(req)
	default:
		err = lberr.NewInvalidArg("req not found protocol type , val is %s", protocolType)
	}
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var headers = make(map[string]string)
	headers[_const.ProtocolType] = _const.PROTO_TYPE_PROTO3

	var target = fmt.Sprintf("%s://%s:%s", "http", node.Host, node.Extra)
	result, err := url.JoinPath(target, path)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	log.Infof("do http request request: %s", req.String())
	resp, err := restysdk.NewRequest().SetHeaders(headers).SetBody(body).Execute(method, result)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	val := resp.Header().Get(_const.ProtocolType)
	switch val {
	case _const.PROTO_TYPE_API_JSON:
		log.Infof("do http resp is %s", string(resp.Body()))
		var respBody Resp
		err := json.Unmarshal(resp.Body(), &respBody)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		if respBody.Errcode > 0 {
			return lberr.NewErr(respBody.Errcode, respBody.Errmsg)
		}
		unmarshaler := &jsonpb.Unmarshaler{AllowUnknownFields: true}
		err = unmarshaler.Unmarshal(bytes.NewReader([]byte(respBody.Data)), out)
	case _const.PROTO_TYPE_PROTO3:
		err = proto.Unmarshal(resp.Body(), out)
		log.Infof("do http resp is %s", out.String())
	default:
		err = lberr.NewInvalidArg("resp not found protocol type , val is %s", val)
	}
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
