package impl

import (
	"context"
	"github.com/oldbai555/bgg/lbblog"
	"testing"
)

func init() {
	InitDbOrm()
}

func TestLbblogServer_AddComment(t *testing.T) {
	lbblogServer.AddComment(context.Background(), &lbblog.AddCommentReq{
		Comment: &lbblog.ModelComment{
			UserId:    4,
			ArticleId: 4,
			Content:   "评论",
		},
	})
}
