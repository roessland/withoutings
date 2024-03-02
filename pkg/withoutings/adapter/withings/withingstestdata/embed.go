package withingstestdata

import _ "embed"

//go:embed withings-resp-sleepv2-getsummary-0.json
var Sleepv2GetsummarySuccess []byte

//go:embed withings-resp-sleepv2-get-0.json
var Sleepv2GetSuccess []byte

//go:embed withings-resp-measure-getmeas-0-weight.json
var MeasureGetmeasSuccess []byte
