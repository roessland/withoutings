package openapi2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type ErrorResponse struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
	Error  string      `json:"error"`
}

// ParseErrorResponse parses an HTTP response from any call
func ParseErrorResponse(rsp *http.Response) error {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)

	defer func() {
		_ = rsp.Body.Close()
		rsp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}()

	if err != nil {
		return err
	}

	response := ErrorResponse{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil
	}

	if response.Status == 0 {
		return nil
	} else {
		return fmt.Errorf("status %d: %s", response.Status, response.Error)
	}
}
