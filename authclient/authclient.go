package authclient

import (
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
	Timeout time.Duration
}

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

//go:generate mockery --name=HTTPClient --filename=http_client_mock.go --disable-version-string
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	HTTPClient HTTPClient
	Logger     Logger
	Cfg        Config
}

func New(cfg Config, logger Logger) *Client {
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
		Transport: otelhttp.NewTransport(&retryRoundTripper{
			next: &rateLimitRoundTripper{
				next: &circuitBreakerRoundTripper{
					next:           http.DefaultTransport,
					circuitBreaker: cb,
				},
				limiter: rate.NewLimiter(rate.Every(500*time.Millisecond), 2),
			},
			maxRetries: 10,
			delay:      100 * time.Millisecond,
		}),
		Timeout: cfg.Timeout,
	}

	return Init(cfg, httpClient, logger)
}

func Init(cfg Config, httpClient HTTPClient, logger Logger) *Client {
	return &Client{
		HTTPClient: httpClient,
		Cfg:        cfg,
		Logger:     logger,
	}
}
