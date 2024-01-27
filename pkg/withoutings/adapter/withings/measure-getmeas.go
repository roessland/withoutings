package withings

import (
	"context"
	"encoding/json"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"io"
	"net/http"
)

// Example when there is no data:
// {"status":0,"body":{"updatetime":1706376751,"timezone":"Europe\/Oslo","measuregrps":[]}}%

// Example with weight data only:
// ❯ curl -X POST -d "action=getmeas&meastypes=1&category=1&startdate=1706375413&&enddate=1706375414" -H "Authorization: Bearer <auth>" https://wbsapi.withings.net/measure
// {"status":0,"body":{"updatetime":1706376751,"timezone":"Europe\/Oslo","measuregrps":[{"grpid":5543254354,"attrib":0,"date":1706375414,"created":1706375434,"modified":1706375434,"category":1,"deviceid":"fdsafdsafdsafdsaf","hash_deviceid":"e8asdfasdfasdfaf3","measures":[{"value":74921,"type":1,"unit":-3,"algo":0,"fm":131}],"modelid":5,"model":"Body+","comment":null}]}}%

func (c *HTTPClient) MeasureGetmeas(ctx context.Context, accessToken string, params withings.MeasureGetmeasParams) (*withings.MeasureGetmeasResponse, error) {
	return c.WithAccessToken(accessToken).MeasureGetmeas(ctx, params)
}

// NewMeasureGetmeasRequest creates a new MeasureGetmeas request.
func (c *HTTPClient) NewMeasureGetmeasRequest(params withings.MeasureGetmeasParams) (*http.Request, error) {
	return c.NewRequest("/measure", nil, []byte((params)))
}

//
// curl command to send req with a body using POST:
// curl -X POST -d "userid=313337&startdate=1706375413&enddate=1706375414&appli=1" https://wbsapi.withings.net/measure?action=getmeas

// MeasureGetmeas gets measures
func (c *AuthenticatedClient) MeasureGetmeas(ctx context.Context, params withings.MeasureGetmeasParams) (*withings.MeasureGetmeasResponse, error) {
	log := logging.MustGetLoggerFromContext(ctx)
	req, err := c.NewMeasureGetmeasRequest(params)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.WithAccessToken(c.AccessToken).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	var resp withings.MeasureGetmeasResponse

	resp.Raw, err = io.ReadAll(httpResp.Body)
	if err != nil {
		log.WithField("event", "error.MeasureGetmeas.readbody.failed").
			WithError(err).
			Error()
		return nil, err
	}

	err = json.Unmarshal(resp.Raw, &resp)
	if err != nil {
		log.WithField("response_body", string(resp.Raw)).
			WithField("event", "error.MeasureGetmeas.unmarshal.failed").
			WithError(err).
			Error()
		return nil, err
	}

	if resp.Status != 0 {
		return &resp, withings.APIError(resp.Status)
	}

	return &resp, nil
}
