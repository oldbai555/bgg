package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Embed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/embeddings", r.URL.Path)

		var req embeddingsRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		assert.Equal(t, "bge-m3", req.Model)
		assert.Equal(t, "hello", req.Prompt)

		_ = json.NewEncoder(w).Encode(embeddingsResponse{Embedding: []float32{0.1, 0.2, 0.3}})
	}))
	defer server.Close()

	client := NewClient(server.URL, "bge-m3", "qwen2.5:7b", 5*time.Second)
	vec, err := client.Embed(context.Background(), "hello")
	require.NoError(t, err)
	assert.Equal(t, []float32{0.1, 0.2, 0.3}, vec)
}

func TestClient_Embed_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(embeddingsResponse{})
	}))
	defer server.Close()

	client := NewClient(server.URL, "bge-m3", "qwen2.5:7b", 5*time.Second)
	_, err := client.Embed(context.Background(), "hello")
	assert.Error(t, err)
}

func TestClient_Chat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/chat", r.URL.Path)

		var req chatRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		assert.Equal(t, "qwen2.5:7b", req.Model)
		assert.False(t, req.Stream)
		require.Len(t, req.Messages, 2)
		assert.Equal(t, "system", req.Messages[0].Role)
		assert.Equal(t, "user", req.Messages[1].Role)

		_ = json.NewEncoder(w).Encode(chatResponse{Message: chatMessage{Role: "assistant", Content: "答案在这里"}})
	}))
	defer server.Close()

	client := NewClient(server.URL, "bge-m3", "qwen2.5:7b", 5*time.Second)
	answer, err := client.Chat(context.Background(), "系统提示", "用户问题")
	require.NoError(t, err)
	assert.Equal(t, "答案在这里", answer)
}

func TestClient_Chat_UpstreamError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("model not found"))
	}))
	defer server.Close()

	client := NewClient(server.URL, "bge-m3", "qwen2.5:7b", 5*time.Second)
	_, err := client.Chat(context.Background(), "系统提示", "用户问题")
	assert.Error(t, err)
}
