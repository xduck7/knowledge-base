package builder

import (
	"net"
	"net/http"
	"time"
)

type Client struct {
	*http.Client
	baseURL string
}

type Builder struct {
	timeout   time.Duration
	baseURL   string
	userAgent string
	maxIdle   int
}

func NewBuilder() *Builder {
	return &Builder{
		timeout: 30 * time.Second,
		maxIdle: 10,
	}
}

func (b *Builder) Timeout(d time.Duration) *Builder {
	b.timeout = d
	return b
}

func (b *Builder) BaseURL(u string) *Builder {
	b.baseURL = u
	return b
}

func (b *Builder) UserAgent(ua string) *Builder {
	b.userAgent = ua
	return b
}

func (b *Builder) MaxIdleConns(n int) *Builder {
	b.maxIdle = n
	return b
}

func (b *Builder) Build() *Client {
	transport := &http.Transport{
		MaxIdleConns:    b.maxIdle,
		IdleConnTimeout: 90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
	}

	c := &Client{
		Client: &http.Client{
			Timeout:   b.timeout,
			Transport: transport,
		},
		baseURL: b.baseURL,
	}

	return c
}

// Example
func Example() {
	client := NewBuilder().
		Timeout(10 * time.Second).
		BaseURL("https://api.example.com").
		UserAgent("my-service/1.0").
		MaxIdleConns(50).
		Build()

	_, _ = client.Get(client.baseURL + "/health")
}
