package templates_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/pkg/web/templates"
	"github.com/stretchr/testify/require"
)

func TestLoadTemplates(t *testing.T) {
	tmpls := templates.NewTemplates(templates.Config{})
	require.NotNil(t, tmpls)

	tmpls = templates.NewTemplates(templates.Config{EmbeddedOnly: true})
	require.NotNil(t, tmpls)
}

func TestRenderEmbeddedAndDisk(t *testing.T) {
	tmplsDisk := templates.NewTemplates(templates.Config{})
	require.NotNil(t, tmplsDisk)
	require.Equal(t, "disk", tmplsDisk.Source())

	tmplsEmbedded := templates.NewTemplates(templates.Config{EmbeddedOnly: true})
	require.NotNil(t, tmplsEmbedded)
	require.Equal(t, "embedded", tmplsEmbedded.Source())

	t.Run("TemplateTest renders in disk mode", func(t *testing.T) {
		buf := &bytes.Buffer{}
		err := tmplsDisk.RenderTemplateTest(context.Background(), buf)
		require.NoError(t, err)
		html := buf.String()
		require.Contains(t, html, "ThisIsTheError")
		require.Contains(t, html, "ThisIsTheTitle")
		require.Contains(t, html, "ThisIsTheContent")
	})

	t.Run("TemplateTest renders in embedded mode", func(t *testing.T) {
		buf := &bytes.Buffer{}
		err := tmplsEmbedded.RenderTemplateTest(context.Background(), buf)
		require.NoError(t, err)
		html := buf.String()
		require.Contains(t, html, "ThisIsTheError")
	})
}

func TestRenderTemplates(t *testing.T) {
	tmpls := templates.NewTemplates(templates.Config{})
	require.NotNil(t, tmpls)

	var buf *bytes.Buffer

	beforeEach := func(_ *testing.T) {
		buf = &bytes.Buffer{}
	}

	t.Run("TemplateTest renders in embedded mode", func(t *testing.T) {
		beforeEach(t)
		err := tmpls.RenderTemplateTest(context.Background(), buf)
		require.NoError(t, err)
	})

	t.Run("Home handles nil vars", func(t *testing.T) {
		beforeEach(t)
		err := tmpls.RenderHomePage(context.Background(), buf, nil)
		require.NoError(t, err)
	})

	t.Run("RenderRefreshAccessToken handles nil vars", func(t *testing.T) {
		beforeEach(t)
		err := tmpls.RenderRefreshAccessToken(context.Background(), buf, nil, "")
		require.NoError(t, err)
	})

	t.Run("RenderRefreshAccessToken shows error", func(t *testing.T) {
		beforeEach(t)
		err := tmpls.RenderRefreshAccessToken(context.Background(), buf, nil, "no worky")
		require.NoError(t, err)
		require.Contains(t, buf.String(), "no worky")
	})

	t.Run("RenderSleepSummaries handles nil vars", func(t *testing.T) {
		beforeEach(t)
		var paramSleepData *sleep.GetSleepSummaryOutput
		paramErr := ""
		err := tmpls.RenderSleepSummaries(context.Background(), buf, paramSleepData, paramErr)
		require.NoError(t, err)
		require.Contains(t, buf.String(), "No data")
	})

	t.Run("RenderSubscriptionsPage handles nil vars", func(t *testing.T) {
		beforeEach(t)
		err := tmpls.RenderSubscriptionsPage(context.Background(), buf, nil, nil, "")
		require.NoError(t, err)
		require.Contains(t, buf.String(), "don't have")
	})

	t.Run("RenderSubscriptionsWithingsPage handles nil vars", func(t *testing.T) {
		beforeEach(t)
		err := tmpls.RenderSubscriptionsWithingsPage(context.Background(), buf, nil, "")
		require.NoError(t, err)
		require.Contains(t, buf.String(), "don't have")
	})

	t.Run("MeasureGetmeasPage handles nil vars", func(t *testing.T) {
		beforeEach(t)
		err := tmpls.RenderMeasureGetmeas(context.Background(), buf, "", "")
		require.NoError(t, err)
		require.Contains(t, buf.String(), "Measure - Getmeas")
	})
}
