// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package knowledge_qa

import (
	"context"
	"fmt"
	"strings"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/internal/vectorstore"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	chunkSize       = 500 // 单个片段目标字数（按 rune 计）
	chunkOverlap    = 50  // 片段间重叠字数，避免关键信息被切断在片段边界
	reindexPageSize = 50
)

type KnowledgeQaReindexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewKnowledgeQaReindexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KnowledgeQaReindexLogic {
	return &KnowledgeQaReindexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// KnowledgeQaReindex 全量重建博客知识库向量索引：拉取全部已发布文章 → 逐篇取全文
// → 分片 → embed → 写入 Redis Vector Set。Phase 1 手动触发，不联动文章发布事件。
func (l *KnowledgeQaReindexLogic) KnowledgeQaReindex(req *types.KnowledgeQaReindexReq) (resp *types.KnowledgeQaReindexResp, err error) {
	index := vectorstore.NewBlogIndex(l.svcCtx.Redis)
	if err := index.Reset(l.ctx); err != nil {
		return nil, err
	}

	var articleCount, chunkCount int64
	for page := int64(1); ; page++ {
		listResp, err := l.svcCtx.ContentRPC.BlogArticleList(l.ctx, &contentclient.BlogArticleListRequest{
			Page:     page,
			PageSize: reindexPageSize,
			Status:   consts.BlogArticleStatusPublished,
		})
		if err != nil {
			return nil, errs.WrapGRPCError("获取已发布文章列表失败", err)
		}
		if len(listResp.List) == 0 {
			break
		}

		for _, item := range listResp.List {
			n, err := l.indexArticle(index, item.Id)
			if err != nil {
				return nil, err
			}
			articleCount++
			chunkCount += n
		}

		if int64(len(listResp.List)) < reindexPageSize {
			break
		}
	}

	return &types.KnowledgeQaReindexResp{
		ArticleCount: articleCount,
		ChunkCount:   chunkCount,
	}, nil
}

func (l *KnowledgeQaReindexLogic) indexArticle(index *vectorstore.BlogIndex, articleID uint64) (int64, error) {
	detail, err := l.svcCtx.ContentRPC.BlogArticleDetail(l.ctx, &contentclient.BlogArticleDetailRequest{Id: articleID})
	if err != nil {
		return 0, errs.WrapGRPCError(fmt.Sprintf("获取文章详情失败（id=%d）", articleID), err)
	}

	chunks := splitIntoChunks(detail.Content, chunkSize, chunkOverlap)
	for i, chunk := range chunks {
		vec, err := l.svcCtx.OllamaClient.Embed(l.ctx, chunk)
		if err != nil {
			return 0, err
		}

		elementID := fmt.Sprintf("%d:%d", articleID, i)
		if err := index.Add(l.ctx, elementID, vec, vectorstore.Attr{
			ArticleID: articleID,
			Title:     detail.Title,
			ChunkText: chunk,
		}); err != nil {
			return 0, err
		}
	}
	return int64(len(chunks)), nil
}

// splitIntoChunks 先按空行分段落，单段超过 size 再按固定长度二次切分并保留 overlap
// 个字的重叠，避免关键信息被整齐切断在片段边界上。
func splitIntoChunks(text string, size, overlap int) []string {
	var chunks []string
	for _, paragraph := range strings.Split(text, "\n\n") {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}
		chunks = append(chunks, splitByRuneLength(paragraph, size, overlap)...)
	}
	return chunks
}

func splitByRuneLength(s string, size, overlap int) []string {
	runes := []rune(s)
	if len(runes) <= size {
		return []string{s}
	}

	var result []string
	step := size - overlap
	for start := 0; start < len(runes); start += step {
		end := start + size
		if end > len(runes) {
			end = len(runes)
		}
		result = append(result, string(runes[start:end]))
		if end == len(runes) {
			break
		}
	}
	return result
}
