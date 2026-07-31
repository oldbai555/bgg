// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package knowledge_qa

import (
	"context"
	"fmt"
	"strings"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/internal/vectorstore"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	askTopK = 5
	// askScoreThreshold 低于此余弦相似度的检索结果视为不相关，不喂给模型——VSIM 总会
	// 返回 topK 个结果，哪怕全都不相关，加这道门槛避免模型拿着不相关片段"言之凿凿"瞎答。
	askScoreThreshold = 0.5

	noRelevantContentAnswer = "抱歉，知识库里没有找到与这个问题相关的内容。"

	systemPrompt = "你是一个只根据给定资料回答问题的助手。只使用下面提供的资料回答用户问题；" +
		"如果资料中没有相关信息，必须如实说明未找到相关内容，禁止编造资料中不存在的内容。"
)

type KnowledgeQaAskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewKnowledgeQaAskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KnowledgeQaAskLogic {
	return &KnowledgeQaAskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// KnowledgeQaAsk：embed 问题 → 向量检索博客知识库 → 过滤低相关度结果 → 命中则拼资料
// 喂给本地模型生成回答，否则直接告知未找到相关内容（不调用模型，避免瞎编）。
func (l *KnowledgeQaAskLogic) KnowledgeQaAsk(req *types.KnowledgeQaAskReq) (resp *types.KnowledgeQaAskResp, err error) {
	question := strings.TrimSpace(req.Question)
	if question == "" {
		return nil, errs.New(errs.CodeBadRequest, "问题不能为空")
	}

	queryVec, err := l.svcCtx.OllamaClient.Embed(l.ctx, question)
	if err != nil {
		return nil, err
	}

	index := vectorstore.NewBlogIndex(l.svcCtx.Redis)
	results, err := index.Search(l.ctx, queryVec, askTopK)
	if err != nil {
		return nil, err
	}

	relevant := filterByScore(results, askScoreThreshold)
	if len(relevant) == 0 {
		return &types.KnowledgeQaAskResp{Answer: noRelevantContentAnswer, Sources: []string{}}, nil
	}

	userPrompt, sources := buildPrompt(question, relevant)

	answer, err := l.svcCtx.OllamaClient.Chat(l.ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}

	return &types.KnowledgeQaAskResp{
		Answer:  answer,
		Sources: sources,
	}, nil
}

func filterByScore(results []vectorstore.SearchResult, threshold float32) []vectorstore.SearchResult {
	relevant := make([]vectorstore.SearchResult, 0, len(results))
	for _, r := range results {
		if r.Score >= threshold {
			relevant = append(relevant, r)
		}
	}
	return relevant
}

// buildPrompt 拼资料上下文 + 问题，同时收集去重后的来源文章标题（供前端展示引用）。
func buildPrompt(question string, results []vectorstore.SearchResult) (prompt string, sources []string) {
	var contextBuilder strings.Builder
	seen := make(map[string]struct{})

	for i, r := range results {
		contextBuilder.WriteString(fmt.Sprintf("【资料%d，来自文章《%s》】\n%s\n\n", i+1, r.Attr.Title, r.Attr.ChunkText))
		if _, ok := seen[r.Attr.Title]; !ok {
			seen[r.Attr.Title] = struct{}{}
			sources = append(sources, r.Attr.Title)
		}
	}

	prompt = fmt.Sprintf("资料：\n%s\n问题：%s", contextBuilder.String(), question)
	return prompt, sources
}
