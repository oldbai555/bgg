package impl

import (
	"github.com/oldbai555/bgg/lbblog"
	"github.com/oldbai555/lbtool/log"
	"testing"
)

func TestReq(t *testing.T) {
	var req = &lbblog.AddArticleReq{
		Article: &lbblog.ModelArticle{
			Title: "test",
		},
	}
	err := req.ValidateAll()
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	return
}
