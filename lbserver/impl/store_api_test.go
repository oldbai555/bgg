package impl

import (
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
	file, err := ConvertMediaFile("image.jpg", "http://mmbiz.qpic.cn/mmbiz_jpg/HFcU55kicfHeKGqRia75aJgmDENTKJGmkbbRK4uwWBhEFrq4fplNoOiaANYlFobrPQqhZqIFnl1r8LkmK1MXw385w/0")
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	log.Infof("url is %s", file)
}
