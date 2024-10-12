package authclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var defaultClient = http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	},
}

type Logger func(ctx context.Context, msg string, v ...any)

type ServiceClient struct {
	url  string
	log  Logger
	http *http.Client
}

func New(url string, log Logger, options ...func(*ServiceClient)) *ServiceClient {
	client := ServiceClient{
		url:  url,
		log:  log,
		http: &defaultClient,
	}

	for _, option := range options {
		option(&client)
	}
	return &client
}

// This function would rather use a custom http(s) client over the default client configured at the package level
func WithCustomClient(httpClient *http.Client) func(*ServiceClient) {
	return func(client *ServiceClient) {
		client.http = httpClient
	}
}

func (cln *ServiceClient) Authenticate(ctx context.Context, authorization string) (AuthenticateResp, error) {

	endpoint := fmt.Sprintf("%s/auth/authenticate", cln.url)
	headers := map[string]string{
		"authorization": authorization,
	}

	var resp AuthenticateResp
	if err := cln.invoke(ctx, http.MethodGet, endpoint, headers, nil, &resp); err != nil {
		return AuthenticateResp{}, err
	}

	return resp, nil
}

func (cln *ServiceClient) invoke(ctx context.Context, method string, url string, headers map[string]string, r io.Reader, v any) error {
	cln.log(ctx, "authclient: invoke: started", "method", method, "url", url)
	defer cln.log(ctx, "authclient: invoke: completed")

	req, err := http.NewRequestWithContext(ctx, method, url, r)
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	for key, value := range headers {
		cln.log(ctx, "authclient: invoke", "key", key, "value", value)
		req.Header.Set(key, value)
	}

	resp, err := cln.http.Do(req)
	if err != nil {
		return fmt.Errorf("do: error: %w", err)
	}
	defer resp.Body.Close()

	cln.log(ctx, "authclient: invoke", "statuscode", resp.StatusCode)

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("copy error: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil

	case http.StatusOK:
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("failed: response: %s, decoding error: %w ", string(data), err)
		}
		return nil

	case http.StatusUnauthorized:
		var err Error
		if err := json.Unmarshal(data, &err); err != nil {
			return fmt.Errorf("failed: response: %s, decoding error: %w ", string(data), err)
		}
		return err

	default:
		return fmt.Errorf("failed: response: %s", string(data))
	}
}
