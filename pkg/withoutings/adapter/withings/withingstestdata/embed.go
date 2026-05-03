package withingstestdata

import _ "embed"

//go:embed withings-resp-sleepv2-getsummary-0.json
var Sleepv2GetsummarySuccess []byte

//go:embed withings-resp-sleepv2-get-0.json
var Sleepv2GetSuccess []byte

//go:embed withings-resp-measure-getmeas-0-weight.json
var MeasureGetmeasSuccess []byte

//go:embed withings-resp-measurev2-getactivity-0.json
var MeasureGetactivitySuccess []byte

//go:embed withings-resp-measurev2-getactivity-0-nodata.json
var MeasureGetactivityNoData []byte

//go:embed withings-resp-measurev2-getintradayactivity-0.json
var MeasureGetintradayactivitySuccess []byte

//go:embed withings-resp-measurev2-getintradayactivity-0-nodata.json
var MeasureGetintradayactivityNoData []byte

//go:embed withings-resp-measurev2-getworkouts-0.json
var MeasureGetworkoutsSuccess []byte

//go:embed withings-resp-measurev2-getworkouts-0-nodata.json
var MeasureGetworkoutsNoData []byte
