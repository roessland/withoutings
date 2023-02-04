package templates_test

import (
	"bytes"
	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/web/templates"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadTemplates(t *testing.T) {
	tmpls := templates.LoadTemplates()
	require.NotNil(t, tmpls)
}

func TestRenderTemplates(t *testing.T) {
	tmpls := templates.LoadTemplates()
	require.NotNil(t, tmpls)

	var buf *bytes.Buffer

	beforeEach := func(t *testing.T) {
		buf = &bytes.Buffer{}
	}

	t.Run("Home handles nil vars", func(t *testing.T) {
		beforeEach(t)
		err := tmpls.RenderHomePage(buf, nil)
		require.NoError(t, err)
	})

	t.Run("RenderRefreshAccessToken handles nil vars", func(t *testing.T) {
		beforeEach(t)
		err := tmpls.RenderRefreshAccessToken(buf, nil, "")
		require.NoError(t, err)
		require.Contains(t, buf.String(), "nil")
	})

	t.Run("RenderRefreshAccessToken shows error", func(t *testing.T) {
		beforeEach(t)
		err := tmpls.RenderRefreshAccessToken(buf, nil, "no worky")
		require.NoError(t, err)
		require.Contains(t, buf.String(), "no worky")
	})

	t.Run("RenderSleepSummaries handles nil vars", func(t *testing.T) {
		beforeEach(t)
		var paramSleepData *sleep.GetSleepSummaryOutput
		var paramErr = ""
		err := tmpls.RenderSleepSummaries(buf, paramSleepData, paramErr)
		require.NoError(t, err)
		require.Contains(t, buf.String(), "No data")
	})

	t.Run("RenderSubscriptionsPage handles nil vars", func(t *testing.T) {
		beforeEach(t)
		err := tmpls.RenderSubscriptionsPage(buf, nil)
		require.NoError(t, err)
		require.Contains(t, buf.String(), "don't have")
	})
}
