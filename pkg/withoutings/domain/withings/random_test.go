package withings_test

import (
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomNonceIsURLSafe(t *testing.T) {
	for i := 0; i < 100; i++ {
		nonce := withings.RandomNonce()
		require.True(t, len(nonce) > 10)
		require.NotContains(t, nonce, "=")
		require.NotContains(t, nonce, "+")
		require.NotContains(t, nonce, "/")
		require.NotContains(t, nonce, " ")
		require.NotContains(t, nonce, "%")
	}
}
