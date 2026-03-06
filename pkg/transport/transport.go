package transport

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"go.uber.org/zap"
)

const defaultUA = "slack-mcp-server/1.0"

// UserAgentTransport wraps another RoundTripper to add User-Agent and cookies
type UserAgentTransport struct {
	roundTripper http.RoundTripper
	userAgent    string
	cookies      []*http.Cookie
	logger       *zap.Logger
}

// NewUserAgentTransport creates a new UserAgentTransport
func NewUserAgentTransport(roundTripper http.RoundTripper, userAgent string, cookies []*http.Cookie, logger *zap.Logger) *UserAgentTransport {
	return &UserAgentTransport{
		roundTripper: roundTripper,
		userAgent:    userAgent,
		cookies:      cookies,
		logger:       logger,
	}
}

// RoundTrip implements the RoundTripper interface
func (t *UserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clonedReq := req.Clone(req.Context())
	clonedReq.Header.Set("User-Agent", t.userAgent)

	for _, cookie := range t.cookies {
		clonedReq.AddCookie(cookie)
	}

	t.logger.Debug("Making request", zap.String("url", clonedReq.URL.String()))

	resp, err := t.roundTripper.RoundTrip(clonedReq)
	if err != nil {
		t.logger.Error("Request failed", zap.Error(err))
	}
	return resp, err
}

// ProvideHTTPClient creates an HTTP client with optional proxy and custom CA support
func ProvideHTTPClient(cookies []*http.Cookie, logger *zap.Logger) *http.Client {
	var proxy func(*http.Request) (*url.URL, error)
	if proxyURL := os.Getenv("SLACK_MCP_PROXY"); proxyURL != "" {
		parsed, err := url.Parse(proxyURL)
		if err != nil {
			logger.Fatal("Failed to parse proxy URL",
				zap.Error(err))
		}
		// Log proxy URL with credentials redacted
		redactedURL := *parsed
		if redactedURL.User != nil {
			redactedURL.User = url.UserPassword("REDACTED", "REDACTED")
		}
		logger.Info("Using proxy", zap.String("proxy_url", redactedURL.String()))
		proxy = http.ProxyURL(parsed)
	}

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	if localCertFile := os.Getenv("SLACK_MCP_SERVER_CA"); localCertFile != "" {
		certs, err := os.ReadFile(localCertFile)
		if err != nil {
			logger.Fatal("Failed to read local certificate file",
				zap.String("cert_file", localCertFile),
				zap.Error(err))
		}
		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			logger.Warn("No certs appended, using system certs only")
		}
	}

	transport := &http.Transport{
		Proxy: proxy,
		TLSClientConfig: &tls.Config{
			RootCAs: rootCAs,
		},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	wrappedTransport := NewUserAgentTransport(transport, defaultUA, cookies, logger)

	client := &http.Client{
		Transport: wrappedTransport,
		Timeout:   30 * time.Second,
	}

	return client
}
