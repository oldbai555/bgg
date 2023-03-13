package impl

import (
	"context"
	"github.com/oldbai555/bgg/client/lbuser"
	webtool2 "github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/lbtool/log"
	"testing"
)

func init() {
	v, err := initViper()
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	lb = &Tool{}
	lb.WebTool, _ = webtool2.NewWebTool(v, webtool2.OptionWithOrm(&lbuser.ModelUser{}), webtool2.OptionWithRdb(), webtool2.OptionWithStorage())
	InitDbOrm()
	val := lb.V.Get("wechatConf")
	err = webtool2.JsonConvertStruct(val, &lb.WechatConf)
	if err != nil {
		log.Errorf("err is %v", err)
	}
}

func TestGetWxGzhMedia(t *testing.T) {
	media, err := GetWxGzhMediaBytes(context.Background(), "")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	log.Infof("resp %v", media)
}
