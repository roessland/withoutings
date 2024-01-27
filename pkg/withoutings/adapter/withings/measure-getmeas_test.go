package withings_test

import (
	withingsHttp "github.com/roessland/withoutings/pkg/withoutings/adapter/withings"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestMeasureGetmeas(t *testing.T) {
	c := withingsHttp.NewClient("a", "b", "c").WithAccessToken("test")
	req, err := c.NewMeasureGetmeasRequest(`userid=313337&startdate=1706375413&enddate=1706375414&appli=1`)
	require.NoError(t, err)
	require.Equal(t, "", req.URL.Query().Get("appli"), "params go in the _body_, not the query string")
	body, _ := io.ReadAll(req.Body)
	require.EqualValues(t, "userid=313337&startdate=1706375413&enddate=1706375414&appli=1", string(body))
}
