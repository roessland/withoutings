package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/roessland/withoutings/withings"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	TypeWithingsAPICall = "withings:apicall"
)

type WithingsAPICallPayload struct {
	Body   []byte
	Header http.Header
	Method string
	Token  withings.Token
	URL    url.URL
}

func NewWithingsAPICallTask(req *http.Request, token withings.Token) (*asynq.Task, error) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(reqBody))

	payload, err := json.Marshal(WithingsAPICallPayload{
		Body:   reqBody,
		Header: req.Header,
		Method: req.Method,
		Token:  withings.Token{},
		URL:    url.URL{},
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(
		TypeWithingsAPICall,
		payload,
		asynq.MaxRetry(5),
		asynq.Timeout(5*time.Minute)), nil
}

type WithingsAPICallProcessor struct {
}

func NewWithingsAPICallProcessor() *WithingsAPICallProcessor {
	return &WithingsAPICallProcessor{}
}

func (processor *WithingsAPICallProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p WithingsAPICallPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Calling but not really %s", p.URL.String())

	// TODO actually call API
	return nil
}
