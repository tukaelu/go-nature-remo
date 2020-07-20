package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	apiBaseURL = "https://api.nature.global/"
	apiVersion = 1
	libVersion = "0.0.1"
	reqTimeout = 30 * time.Second
)

// Client for Nature Remo API.
type Client struct {
	HTTPClient      *http.Client
	BaseURL         string
	AccessToken     string
	UserAgent       string
	LatestRateLimit *RateLimit
	Version         string
}

// RateLimit is the request limit expressed in the X-Rate-Limit-* headers.
// See https://developer.nature.global/
type RateLimit struct {
	Limit     int64
	Reset     time.Time
	Remaining int64
}

// NewClient creates a client for the NatureRemo API.
func NewClient(accessToken string) *Client {
	return &Client{
		BaseURL:     fmt.Sprintf("%s%d", apiBaseURL, apiVersion),
		AccessToken: accessToken,
		UserAgent:   fmt.Sprintf("tukaelu/go-nature-remo (Ver: %s)", libVersion),
		Version:     libVersion,
		HTTPClient:  &http.Client{},
	}
}

// Get is an implementation of the HTTP GET method.
func (cli *Client) Get(ctx context.Context, path string, params url.Values, p interface{}) error {
	endpoint := fmt.Sprintf("%s/%s", cli.BaseURL, path)
	if params != nil {
		endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	res, err := cli.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rateLimit, err := cli.parseLimitHeader(res.Header)
	if err != nil {
		return err
	}
	cli.LatestRateLimit = rateLimit

	if !(res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices) {
		reason, err := ioutil.ReadAll(res.Body)
		if err != nil || len(reason) == 0 {
			return fmt.Errorf("Request failed: Status=%d (no reason)", res.StatusCode)
		}
		return fmt.Errorf("Request failed: Status=%d, Error= %s", res.StatusCode, string(reason))
	}

	if err := json.NewDecoder(res.Body).Decode(p); err != nil {
		return fmt.Errorf("Failed to parse the response. (%s)", err.Error())
	}

	return nil
}

func (cli *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.AccessToken))
	req.Header.Set("User-Agent", cli.UserAgent)
	cli.HTTPClient.Timeout = reqTimeout
	return cli.HTTPClient.Do(req)
}

func (cli *Client) parseLimitHeader(h http.Header) (*RateLimit, error) {
	limit := h.Get("X-Rate-Limit-Limit")
	if limit == "" {
		return nil, fmt.Errorf("X-Rate-Limit-Limit header was not responded")
	}
	vLimit, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("X-Rate-Limit-Limit is invalid: %s", limit)
	}

	reset := h.Get("X-Rate-Limit-Reset")
	if reset == "" {
		return nil, fmt.Errorf("X-Rate-Limit-Reset header was not responded")
	}
	vReset, err := strconv.ParseInt(reset, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("X-Rate-Limit-Reset is invalid: %s", reset)
	}

	remaining := h.Get("X-Rate-Limit-Remaining")
	if remaining == "" {
		return nil, fmt.Errorf("X-Rate-Limit-Remaining header was not responded")
	}
	vRemaining, err := strconv.ParseInt(remaining, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("X-Rate-Limit-Remaining is invalid: %s", remaining)
	}

	return &RateLimit{
		Limit:     vLimit,
		Reset:     time.Unix(vReset, 0),
		Remaining: vRemaining,
	}, nil
}
