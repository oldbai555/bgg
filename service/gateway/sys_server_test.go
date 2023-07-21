package main

import (
	"github.com/oldbai555/bgg/service/lb"
	"github.com/oldbai555/bgg/service/lbblog"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/restysdk"
	"testing"
)

func Test_revProxy(t *testing.T) {
	resp, err := restysdk.NewRequest().SetBody(&lbblog.GetArticleListReq{
		Options: &lb.Options{
			OptList: []*lb.Options_Opt{},
			Size:    1,
			Page:    1,
		},
	}).Post("http://localhost:20000/gateway/lbblog/GetArticleList")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	log.Infof("%s", resp.Body())
}
