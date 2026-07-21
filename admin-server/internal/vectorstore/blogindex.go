// Package vectorstore 封装本机 Redis 8 原生 Vector Set（VADD/VSIM/DEL）读写博客知识库
// 向量索引，供 AI 知识库问答（ai/knowledge_qa）使用。详见 admin-server/docs/ai-knowledge-qa-spec.md。
//
// 本机 Redis 8.8.0 只加载了 vectorset 模块，没有 RediSearch，所以走 VADD/VSIM 原生命令，
// 不是 FT.CREATE/FT.SEARCH 那一套；项目锁定的 go-zero v1.9.3 的 *redis.Redis 没有通用的
// Do/DoCtx 方法（v1.10+ 才有），只能借 EvalCtx 跑一段 `redis.call(cmd, KEYS[1], unpack(ARGV))`
// 的 Lua 脚本转发命令——VADD/VSIM 已经在本机验证过可以被 EVAL 调用（未标记 noscript）。
package vectorstore

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

const blogIndexKey = "ai:kb:blog"

// Attr 是每个向量元素挂载的业务属性（Vector Set 的 SETATTR），用来省掉一张 MySQL 表。
type Attr struct {
	ArticleID uint64 `json:"articleId"`
	Title     string `json:"title"`
	ChunkText string `json:"chunkText"`
}

type SearchResult struct {
	ElementID string
	Score     float32
	Attr      Attr
}

type BlogIndex struct {
	rdb *redis.Redis
}

func NewBlogIndex(rdb *redis.Redis) *BlogIndex {
	return &BlogIndex{rdb: rdb}
}

// evalCommand 用一段透传脚本把 cmd 连同 args 转发给 redis.call，规避 go-zero v1.9.3
// *redis.Redis 没有通用 Do/DoCtx 方法的限制。
func (idx *BlogIndex) evalCommand(ctx context.Context, cmd string, args ...any) (any, error) {
	script := "return redis.call('" + cmd + "', KEYS[1], unpack(ARGV))"
	return idx.rdb.EvalCtx(ctx, script, []string{blogIndexKey}, args...)
}

// Reset 清空整个知识库索引，Reindex 流程每次全量重建前调用，避免旧文章/旧分片残留。
func (idx *BlogIndex) Reset(ctx context.Context) error {
	if _, err := idx.rdb.DelCtx(ctx, blogIndexKey); err != nil {
		return errs.Wrap(errs.CodeBadDB, "清空向量索引失败", err)
	}
	return nil
}

// Add 写入一个向量片段。
func (idx *BlogIndex) Add(ctx context.Context, elementID string, vec []float32, attr Attr) error {
	attrJSON, err := json.Marshal(attr)
	if err != nil {
		return errs.Wrap(errs.CodeInternalError, "向量属性序列化失败", err)
	}

	args := make([]any, 0, 4+len(vec))
	args = append(args, "VALUES", len(vec))
	for _, v := range vec {
		args = append(args, v)
	}
	args = append(args, elementID, "SETATTR", string(attrJSON))

	if _, err := idx.evalCommand(ctx, "VADD", args...); err != nil {
		return errs.Wrap(errs.CodeBadDB, "写入向量索引失败", err)
	}
	return nil
}

// Search 按余弦相似度检索最相关的 topK 个片段（VSIM 默认即余弦相似度）。
func (idx *BlogIndex) Search(ctx context.Context, queryVec []float32, topK int) ([]SearchResult, error) {
	args := make([]any, 0, 6+len(queryVec))
	args = append(args, "VALUES", len(queryVec))
	for _, v := range queryVec {
		args = append(args, v)
	}
	args = append(args, "COUNT", topK, "WITHSCORES", "WITHATTRIBS")

	res, err := idx.evalCommand(ctx, "VSIM", args...)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "向量检索失败", err)
	}
	return parseSearchResult(res)
}

func parseSearchResult(res any) ([]SearchResult, error) {
	if res == nil {
		return nil, nil
	}
	arr, ok := res.([]any)
	if !ok {
		return nil, errs.New(errs.CodeBadDB, fmt.Sprintf("向量检索返回类型异常: %T", res))
	}
	if len(arr)%3 != 0 {
		return nil, errs.New(errs.CodeBadDB, "向量检索返回字段数异常（预期 element/score/attrib 三元组）")
	}

	results := make([]SearchResult, 0, len(arr)/3)
	for i := 0; i+2 < len(arr); i += 3 {
		elementID, _ := arr[i].(string)

		score, err := toFloat32(arr[i+1])
		if err != nil {
			return nil, errs.Wrap(errs.CodeBadDB, "向量检索分数解析失败", err)
		}

		var attr Attr
		if attrRaw, _ := arr[i+2].(string); attrRaw != "" {
			if err := json.Unmarshal([]byte(attrRaw), &attr); err != nil {
				return nil, errs.Wrap(errs.CodeBadDB, "向量检索属性解析失败", err)
			}
		}

		results = append(results, SearchResult{ElementID: elementID, Score: score, Attr: attr})
	}
	return results, nil
}

func toFloat32(v any) (float32, error) {
	switch t := v.(type) {
	case string:
		f, err := strconv.ParseFloat(t, 32)
		if err != nil {
			return 0, err
		}
		return float32(f), nil
	case float64:
		return float32(t), nil
	case int64:
		return float32(t), nil
	default:
		return 0, fmt.Errorf("未知的分数类型 %T", v)
	}
}
