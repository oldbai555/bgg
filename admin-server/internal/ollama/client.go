// Package ollama 封装本地 Ollama REST API 的最小客户端，供 AI 知识库问答（ai/knowledge_qa）
// 使用：Embed 做文本向量化，Chat 做 RAG 生成。详见 admin-server/docs/ai-knowledge-qa-spec.md。
package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"postapocgame/admin-server/pkg/errs"
)

type Client struct {
	baseURL    string
	embedModel string
	chatModel  string
	httpClient *http.Client
}

func NewClient(baseURL, embedModel, chatModel string, timeout time.Duration) *Client {
	return &Client{
		baseURL:    baseURL,
		embedModel: embedModel,
		chatModel:  chatModel,
		httpClient: &http.Client{Timeout: timeout},
	}
}

type embeddingsRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type embeddingsResponse struct {
	Embedding []float32 `json:"embedding"`
}

// Embed 调用 /api/embeddings 把文本转成向量。
func (c *Client) Embed(ctx context.Context, text string) ([]float32, error) {
	reqBody, err := json.Marshal(embeddingsRequest{Model: c.embedModel, Prompt: text})
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "embedding 请求体序列化失败", err)
	}

	resp, err := c.post(ctx, "/api/embeddings", reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out embeddingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, errs.Wrap(errs.CodeBadGateway, "ollama embedding 响应解析失败", err)
	}
	if len(out.Embedding) == 0 {
		return nil, errs.New(errs.CodeBadGateway, "ollama embedding 返回为空")
	}
	return out.Embedding, nil
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type chatResponse struct {
	Message chatMessage `json:"message"`
}

// Chat 调用 /api/chat 做非流式对话生成，返回助手回复的纯文本。
func (c *Client) Chat(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	reqBody, err := json.Marshal(chatRequest{
		Model: c.chatModel,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Stream: false,
	})
	if err != nil {
		return "", errs.Wrap(errs.CodeInternalError, "chat 请求体序列化失败", err)
	}

	resp, err := c.post(ctx, "/api/chat", reqBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", errs.Wrap(errs.CodeBadGateway, "ollama chat 响应解析失败", err)
	}
	if out.Message.Content == "" {
		return "", errs.New(errs.CodeBadGateway, "ollama chat 返回为空")
	}
	return out.Message.Content, nil
}

func (c *Client) post(ctx context.Context, path string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "构造 ollama 请求失败", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadGateway, "AI 服务暂不可用，请稍后重试", err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		msg, _ := io.ReadAll(resp.Body)
		return nil, errs.New(errs.CodeBadGateway, fmt.Sprintf("ollama 返回异常状态码 %d: %s", resp.StatusCode, string(msg)))
	}
	return resp, nil
}
