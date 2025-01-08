/**
 * @Author: zjj
 * @Date: 2024/12/31
 * @Desc:
**/

package lbsingleserver

import (
	"context"
	"github.com/oldbai555/bgg/pkg/bctx"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/bgg/service/lbsingleserver/cache"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/json"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/pkg/restysdk"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/uctx"
	"testing"
)

func init() {
	cache.InitCache()
	Init()
}

type Resp struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
	Data   struct {
		Recommend []struct {
			Name             string   `json:"name"`
			LiteratureAuthor string   `json:"literature_author"`
			Sid              string   `json:"sid"`
			Type             string   `json:"type"`
			Body             []string `json:"body"`
			Tag              []string `json:"tag"`
			Img              string   `json:"img"`
			LikeCount        int      `json:"like_count"`
			VocabCount       int      `json:"vocab_count"`
		} `json:"recommend"`
	} `json:"data"`
}

func TestInit(t *testing.T) {
	request := restysdk.NewRequest()
	response, err := request.SetQueryParams(map[string]string{
		"query":      "每日金句",
		"src_id":     "51328",
		"query_type": "exact",
		"type":       "sentence",
		"pn":         "1",
		"ps":         "1000",
		"smpid":      "",
		"tab_type":   "",
		"gssda_res":  "{}",
	}).Get("https://hanyu.baidu.com/hanyu/api/sentencelistv2")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	var resp Resp
	err = json.Unmarshal(response.Body(), &resp)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	ctx := bctx.NewCtx(context.Background())
	for _, val := range resp.Data.Recommend {
		OrmDailyShortSentences.NewBaseScope().UpdateOrCreate(ctx, map[string]interface{}{
			lbsingle.FieldContent_: val.Body[0],
			lbsingle.FieldType_:    int32(lbsingle.ModelDailyShortSentences_TypeQuote),
		}, map[string]interface{}{
			lbsingle.FieldImg_:              val.Img,
			lbsingle.FieldLiteratureAuthor_: val.LiteratureAuthor,
		})
	}
}

func TestTransaction(t *testing.T) {
	ctx := bctx.NewCtx(context.Background())
	err := gormx.OnTransaction(ctx, func(ctx uctx.IUCtx, trId string) error {
		err := OrmDailyShortSentences.WithTransactionId(trId).Create(ctx, &lbsingle.ModelDailyShortSentences{
			Content: "测试事务的短语2",
		})
		if err != nil {
			t.Errorf("err:%v", err)
			return err
		}
		return lberr.NewInvalidArg("%s failed", trId)
	})
	if err != nil {
		t.Errorf("err:%v", err)
	}
	err = OrmDailyShortSentences.Create(ctx, &lbsingle.ModelDailyShortSentences{
		Content: "测试事务的短语3",
	})
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}
