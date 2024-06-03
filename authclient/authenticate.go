package authclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/opentracing/opentracing-go"

	"net/http"
)

func (ac *Client) Authenticate(ctx context.Context, token string) (user *CurrentUser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "chatservers.Authenticate")
	defer span.Finish()

	data := map[string]string{"token": token}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := jaegertracing.NewTracedRequest(http.MethodPost, ac.Cfg.AuthURL, bytes.NewBuffer(jsonData), span)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Log the request headers for debugging purposes
	for k, v := range req.Header {
		ac.Logger.Infof("[AUTH REQUEST HTTP HEADER] %s: %s\n", k, v)
	}

	resp, err := ac.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	span.SetTag("user_id", user.ID)

	return user, nil
}
