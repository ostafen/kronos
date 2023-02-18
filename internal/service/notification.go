package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type NotificationService interface {
	Send(ctx context.Context, url string, payload any) error
}

type httpNotificationService struct {
}

func NewNotificationService() NotificationService {
	return &httpNotificationService{}
}

func isSuccess(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func (s *httpNotificationService) Send(ctx context.Context, url string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if !isSuccess(resp) {
		return fmt.Errorf("webhook notification to %s failed with status: %s", url, resp.Status)
	}
	return nil
}
