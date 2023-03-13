package impl

import (
	"context"
	"github.com/oldbai555/bgg/client/lbuser"
	"github.com/oldbai555/bgg/pkg/webtool"
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
	lb.WebTool, _ = webtool.NewWebTool(v, webtool.OptionWithOrm(&lbuser.ModelUser{}), webtool.OptionWithRdb(), webtool.OptionWithStorage())
	InitDbOrm()
}

func TestConvertMediaFile(t *testing.T) {
	// file, err := ConvertMediaUrl("image.jpg", "http://mmbiz.qpic.cn/mmbiz_jpg/HFcU55kicfHeKGqRia75aJgmDENTKJGmkbbRK4uwWBhEFrq4fplNoOiaANYlFobrPQqhZqIFnl1r8LkmK1MXw385w/0")
	// if err != nil {
	// 	log.Errorf("err is %v", err)
	// 	return
	// }
	// log.Infof("url is %s", file)

	bytes, err := GetWxGzhMediaBytes(context.TODO(), "BlhPDwckoSwoe--2pAsT5a6PJ9-ClgthHyn6h-r_qcG2adFJdnjPrXR6qlRsv2tm")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	url, err := ConvertMediaBytes("1.amr", bytes)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	log.Infof("url is %v", url)

}
