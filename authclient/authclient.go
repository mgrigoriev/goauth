package authclient

import (
	"github.com/sony/gobreaker"
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
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "HTTP Client",
		MaxRequests: 1,
		Interval:    time.Duration(60) * time.Second,
		Timeout:     time.Duration(30) * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
	})

	httpClient := &http.Client{
		Transport: &retryRoundTripper{
			next: &rateLimitRoundTripper{
				next: &circuitBreakerRoundTripper{
					next:           http.DefaultTransport,
					circuitBreaker: cb,
				},
				limiter: rate.NewLimiter(rate.Every(500*time.Millisecond), 2),
			},
			maxRetries: 10,
			delay:      100 * time.Millisecond,
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
