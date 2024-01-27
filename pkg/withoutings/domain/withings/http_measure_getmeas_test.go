package withings_test

import (
	"encoding/json"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"testing"
)

func TestWat(t *testing.T) {
	const resp = `{"status":0,"body":{"updatetime":1706379233,"timezone":"Europe\/Oslo","measuregrps":[{"grpid":5207353660,"attrib":0,"date":1706379213,"created":1706379233,"modified":1706379233,"category":1,"deviceid":"e81","hash_deviceid":"e81","measures":[{"value":75019,"type":1,"unit":-3,"algo":0,"fm":131}],"modelid":5,"model":"Body+","comment":null}]}}`
	var r withings.MeasureGetmeasResponse
	err := json.Unmarshal([]byte(resp), &r)
	if err != nil {
		t.Fatal(err)
	}
}
