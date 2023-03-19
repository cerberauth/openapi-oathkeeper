package generator

import (
	"context"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func TestGenerateFromOpenAPI(t *testing.T) {
	doc, err := openapi3.NewLoader().LoadFromFile(path.Join(basepath, "../test/stub/petstore.openapi.json"))
	require.NoError(t, err)

	ctx := context.Background()
	rules, err := New().Document(doc).Generate(ctx)
	require.NoError(t, err)
	assert.Equal(t, string(rules), "[]")
}
