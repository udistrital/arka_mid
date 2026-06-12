package request

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/udistrital/arka_mid/utils_oas/xray"
)

const (
	authorizationKey = "Authorization"
	contentTypeKey   = "Content-Type"
	acceptHeader     = "Accept"
	contentTypeJSON  = "application/json"
)

var defClient = &http.Client{}
var defaultClient = &http.Client{Timeout: 30 * time.Second}

var ErrResponseDecode = errors.New("response body could not be decoded into target")

// GetWithContext makes a GET request to the given URL using the provided context.
// Checks for non-2xx HTTP status codes, and decodes the response body into target.
func GetWithContext(ctx context.Context, urlp string, target any) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlp, nil)
	if err != nil {
		return 0, fmt.Errorf("could not create request: %w", err)
	}

	resp, err := doRequest(defaultClient, req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusIMUsed {
		return resp.StatusCode, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return resp.StatusCode, fmt.Errorf("%w: %w", ErrResponseDecode, err)
	}

	return resp.StatusCode, nil
}

func SendJson(url string, method string, target, body any) error {
	b := new(bytes.Buffer)
	if body != nil {
		json.NewEncoder(b).Encode(body)
	}

	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set(acceptHeader, contentTypeJSON)
	req.Header.Set(contentTypeKey, contentTypeJSON)

	resp, err := execRequest(defClient, req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	return decodeJSONResponse(resp, target)
}

func GetJsonWSO2(urlp string, target interface{}) error {
	b := new(bytes.Buffer)
	//proxyUrl, err := url.Parse("http://10.20.4.15:3128")
	//http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlp, b)
	if err != nil {
		logs.Error("Error creating request. ", err)
		return err
	}

	req.Header.Set("Accept", "application/json")
	r, err := client.Do(req)
	//r, err := http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		logs.Error("error", err)
		return err
	}
	defer r.Body.Close()

	return decodeJSONResponse(r, target)
}

func GetJson(url string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	resp, err := execRequest(defClient, req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	return decodeJSONResponse(resp, target)
}

func GetJsonTest(url string, target any) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := execRequest(client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, json.NewDecoder(resp.Body).Decode(target)
}

func decodeJSONResponse(resp *http.Response, target interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("http %d: %s", resp.StatusCode, resumirBodyError(body))
	}

	if target == nil || len(bytes.TrimSpace(body)) == 0 {
		return nil
	}

	return json.Unmarshal(body, target)
}

func resumirBodyError(body []byte) string {
	msg := strings.TrimSpace(string(body))
	if msg == "" {
		return "respuesta vacía"
	}

	if len(msg) > 1024 {
		return msg[:1024]
	}

	return msg
}

// doRequest executes req using the provided HTTP client, wrapping the call with
// an X-Ray subsegment scoped to the request's context. The caller is
// responsible for closing resp.Body.
// If the context carries an Authorization value (via WithAuthorization), it is forwarded.
func doRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	req.Header.Set(acceptHeader, contentTypeJSON)
	if token, ok := ctx.Value(authorizationKey).(string); ok && token != "" {
		req.Header.Set(authorizationKey, token)
	}

	ctx, subseg := xray.BeginSubsegment(ctx, req)
	resp, err := client.Do(req.WithContext(ctx))
	xray.CloseSubsegment(subseg, resp, err)

	return resp, err
}

// execRequest executes req using the provided HTTP client, wrapping the call
// with an X-Ray subsegment via the local xray package.
// The caller is responsible for closing resp.Body on success.
func execRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	ctx, subseg := xray.BeginSegmentSec(req)
	resp, err := client.Do(req.WithContext(ctx))
	xray.CloseSubsegment(subseg, resp, err)
	return resp, err
}
