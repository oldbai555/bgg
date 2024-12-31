/**
 * @Author: zjj
 * @Date: 2024/12/31
 * @Desc:
**/

package lbsingleserver

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/pkg/bctx"
	"github.com/oldbai555/bgg/pkg/constant"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/bgg/service/lbsingleserver/cache"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/json"
	"github.com/oldbai555/lbtool/pkg/restysdk"
	"github.com/oldbai555/lbtool/utils"
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
		RetArray []struct {
			List []struct {
				Name             string        `json:"name"`
				LiteratureAuthor string        `json:"literature_author"`
				Sid              string        `json:"sid"`
				Tag              []interface{} `json:"tag"`
				Type             string        `json:"type"`
				Body             []string      `json:"body"`
				IsLike           int           `json:"is_like"`
				IsVocab          int           `json:"is_vocab"`
				LikeCount        int           `json:"like_count"`
				VocabCount       int           `json:"vocab_count"`
			} `json:"list"`
			Type      string        `json:"type"`
			Name      string        `json:"name"`
			HighLight string        `json:"high_light"`
			Filters   []interface{} `json:"filters"`
			QueryType interface{}   `json:"queryType"`
			Count     int           `json:"count"`
		} `json:"ret_array"`
		SrcId     int    `json:"src_id"`
		QueryType string `json:"query_type"`
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
		"ps":         "1",
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
	OrmDailyShortSentences.NewBaseScope().Where(fmt.Sprintf("%s != ''", lbsingle.FieldImg_)).Chunk(ctx, 2000, func(out []*lbsingle.ModelDailyShortSentences) error {
		for _, dailyShortSentences := range out {
			sortUrl, err := downloadFile(dailyShortSentences.Img, constant.BaseStoragePath, utils.GenRandomStr()+".jpg")
			if err != nil {
				log.Errorf("err:%v,%s", err, dailyShortSentences.Img)
				return err
			}
			OrmDailyShortSentences.NewBaseScope().Where(lbsingle.FieldId_, dailyShortSentences.Id).Update(ctx, map[string]interface{}{
				lbsingle.FieldImg_: "https://oldbai.top/oss/download/" + sortUrl,
			})
		}
		return nil
	})
}
