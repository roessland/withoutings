package withingsapi_test

import (
	"github.com/roessland/withoutings/pkg/withoutings/clients/withingsapi"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomNonceIsURLSafe(t *testing.T) {
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
