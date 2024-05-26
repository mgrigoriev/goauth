package authclient

import (
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

type CurrentUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Config struct {
	AuthURL string
}

//go:generate mockery --name=HTTPClient --filename=http_client_mock.go --disable-version-string
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	HTTPClient HTTPClient
	Cfg        Config
}

func New(cfg Config) *Client {
	httpClient := &http.Client{
		Transport: &retryRoundTripper{
			next: &rateLimitRoundTripper{
				next:    http.DefaultTransport,
				limiter: rate.NewLimiter(rate.Every(100*time.Millisecond), 2),
			},
			maxRetries: 10,
			delay:      10 * time.Millisecond,
		},
	}
	return Init(cfg, httpClient)
}

func Init(cfg Config, httpClient HTTPClient) *Client {
	return &Client{
		HTTPClient: httpClient,
		Cfg:        cfg,
	}
}
