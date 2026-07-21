package vectorstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSearchResult(t *testing.T) {
	res := []any{
		"1:0", "0.9969418644905090", `{"articleId":1,"title":"标题A","chunkText":"片段A"}`,
		"2:0", "0.5547809898853302", `{"articleId":2,"title":"标题B","chunkText":"片段B"}`,
	}

	results, err := parseSearchResult(res)
	require.NoError(t, err)
	require.Len(t, results, 2)

	assert.Equal(t, "1:0", results[0].ElementID)
	assert.InDelta(t, 0.9969, results[0].Score, 1e-4)
	assert.Equal(t, uint64(1), results[0].Attr.ArticleID)
	assert.Equal(t, "标题A", results[0].Attr.Title)

	assert.Equal(t, "2:0", results[1].ElementID)
	assert.Equal(t, uint64(2), results[1].Attr.ArticleID)
}

func TestParseSearchResult_Empty(t *testing.T) {
	results, err := parseSearchResult(nil)
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestParseSearchResult_MalformedLength(t *testing.T) {
	_, err := parseSearchResult([]any{"1:0", "0.9"})
	assert.Error(t, err)
}

func TestToFloat32(t *testing.T) {
	cases := []struct {
		name  string
		input any
		want  float32
	}{
		{"string", "0.5", 0.5},
		{"float64", float64(0.75), 0.75},
		{"int64", int64(1), 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := toFloat32(c.input)
			require.NoError(t, err)
			assert.Equal(t, c.want, got)
		})
	}
}

func TestToFloat32_Unknown(t *testing.T) {
	_, err := toFloat32(struct{}{})
	assert.Error(t, err)
}
