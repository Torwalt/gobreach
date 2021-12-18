package hibp

import (
	"encoding/json"
	"fmt"
	"gobreach/internal/domains/breach"
	"net/http"
	"net/url"
	"time"
)

// placeholder is an account, e.g. an email address
const accountBreachEndpoint = "breachedaccount/%s?truncateResponse=false"

type HttpGetter interface {
	Do(req *http.Request) (*http.Response, error)
}

type HIBPClient struct {
	httpGetter HttpGetter
	config     hibpConfig
}

func NewClient(httpGetter HttpGetter, config hibpConfig) *HIBPClient {
	return &HIBPClient{
		httpGetter: httpGetter,
		config:     config,
	}
}

func (c *HIBPClient) GetByEmail(email string) (*breach.Breach, error) {
	ep := fmt.Sprintf(accountBreachEndpoint, email)
	url := &url.URL{Scheme: "https", Host: c.config.BaseURL, Path: ep}

	req := &http.Request{Method: http.MethodGet, URL: url}
	resp, err := c.httpGetter.Do(req)

	if err != nil {
		return nil, breach.NewErrorf(breach.ErrorCodeDatasourceError, "error getting breaches for %s: %s", email, err)
	}

	hrb := &HIBPResponseBody{}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(hrb)

	b, err := c.toBreach(hrb, email)
	if err != nil {
		return nil, breach.NewErrorf(breach.ErrorCodeDatasourceError, "error transforming to Breach %s: %s", email, err)
	}
	return b, nil
}

func (c *HIBPClient) toBreach(hrb *HIBPResponseBody, email string) (*breach.Breach, error) {
	bd, err := time.Parse(time.RFC3339, hrb.BreachDate)
	if err != nil {
		return nil, err
	}
	b := breach.New(
		email,
		hrb.Domain,
		hrb.DataClasses,
		bd,
		hrb.Title,
	)
	return b, nil
}
