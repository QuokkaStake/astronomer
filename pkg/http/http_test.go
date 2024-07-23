package http

import (
	loggerPkg "main/pkg/logger"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpClientErrorCreating(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	client := NewClient(logger, "chain")
	queryInfo, err := client.Get("://test", nil)
	require.Error(t, err)
	require.False(t, queryInfo.Success)
}
