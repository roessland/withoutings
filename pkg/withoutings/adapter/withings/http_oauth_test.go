package withings

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRefreshTokenUnmarshal(t *testing.T) {
	resp := []byte(`{
		"status":0,
		"body":{
		"userid":133337,
		"access_token":"b123456789a",
		"refresh_token":"3425678965432145367f",
		"scope":"user.info,user.activity,user.metrics,user.sleepevents",
		"expires_in":10800,
		"token_type":"Bearer"}
	}`)
	refreshTokenResponse := getRefreshTokenResponse{}
	err := json.Unmarshal(resp, &refreshTokenResponse)
	require.NoError(t, err)
	require.Equal(t, 133337, refreshTokenResponse.Body.UserID)
}
