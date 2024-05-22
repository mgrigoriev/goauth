package authclient

import (
	"io"
	"net/http"
)

type Config struct {
	AuthURL string
}

//go:generate mockery --name=HTTPClient --filename=http_client_mock.go --disable-version-string
type HTTPClient interface {
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

type Client struct {
	HTTPClient HTTPClient
	Cfg        Config
}

func New(cfg Config, httpClient HTTPClient) *Client {
	return &Client{
		HTTPClient: httpClient,
		Cfg:        cfg,
	}
}
