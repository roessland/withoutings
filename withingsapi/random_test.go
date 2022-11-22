package withingsapi_test

import (
	"github.com/roessland/withoutings/withingsapi"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomNonce(t *testing.T) {
	for i := 0; i < 100; i++ {
		nonce := withingsapi.RandomNonce()
		require.True(t, len(nonce) > 10)
		require.NotContains(t, nonce, "=")
		require.NotContains(t, nonce, "+")
		require.NotContains(t, nonce, "/")
		require.NotContains(t, nonce, " ")
		require.NotContains(t, nonce, "%")
	}
}
