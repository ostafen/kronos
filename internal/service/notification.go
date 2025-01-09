package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type NotificationService interface {
	Send(ctx context.Context, url string, payload any) (int, error)
}

type httpNotificationService struct {
}

func NewNotificationService() NotificationService {
	return &httpNotificationService{}
}

func isSuccess(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func (s *httpNotificationService) Send(ctx context.Context, url string, payload any) (int, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data)) // TODO: add support for other types of methods
	if err != nil {
		return -1, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	if !isSuccess(resp) {
		err = fmt.Errorf("webhook notification to %s failed with status: %s", url, resp.Status)
	}
	return resp.StatusCode, err
}
