package hibp

import (
	"encoding/json"
	"fmt"
	"gobreach/internal/domains/breach"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

// placeholder is an account, e.g. an email address
const accountBreachEndpoint = "/breachedaccount/%s"
const apiPath = "/api/v3"

type HttpGetter interface {
	Do(req *http.Request) (*http.Response, error)
}

type sleeper func(d time.Duration)

type hibpClient struct {
	httpGetter HttpGetter
	config     hibpConfig
	retries    int
	sleeper    sleeper
}

func NewClient(httpGetter HttpGetter, config hibpConfig, s sleeper) *hibpClient {
	return &hibpClient{
		httpGetter: httpGetter,
		config:     config,
		retries:    0,
		sleeper:    s,
	}
}

// Retrieve all Breaches for the given email address from HIBP.
// Retry the request 'hibpConfig.maxRetries' times and return an error if there is still a 429.
// This is expected to be run non-concurrent, as there are strict request limits for the API.
func (c *hibpClient) GetByEmail(email string) ([]breach.Breach, *breach.Error) {
	ep := fmt.Sprintf(accountBreachEndpoint, email)
	fullPath := apiPath + ep
	URL := &url.URL{Scheme: "https", Host: c.config.host, Path: fullPath, RawQuery: "truncateResponse=false"}

	req, _ := http.NewRequest(http.MethodGet, URL.String(), nil)
	req.Header.Add("hibp-api-key", c.config.apiKey)
	req.Header.Add("user-agent", "gobreach")

	bS, err := c.getByEmail(req, email)

	if err != nil {
		return nil, err
	}

	return bS, nil
}

func (c *hibpClient) getByEmail(req *http.Request, email string) ([]breach.Breach, *breach.Error) {
	resp, err := c.makeRequest(req, email)
	if err != nil {
		return nil, err
	}
	bS := []breach.Breach{}
	hrbS := &[]hibp200ResponseBreach{}

	defer resp.Body.Close()

	errr := json.NewDecoder(resp.Body).Decode(hrbS)
	if errr != nil {
		return nil, breach.NewErrorf(breach.DatasourceErr, "error decoding breaches for %s: %s", email, err)
	}

	for _, h := range *hrbS {
		b, err := c.toBreach(&h, email)
		if err != nil {
			return nil, breach.NewErrorf(breach.DatasourceErr, "error transforming to Breach %s: %s", email, err)
		}
		bS = append(bS, *b)
	}

	return bS, nil
}

func (c *hibpClient) makeRequest(req *http.Request, email string) (*http.Response, *breach.Error) {
	resp, err := c.httpGetter.Do(req)
	if err != nil {
		return nil, breach.NewErrorf(
			breach.DatasourceErr, "error getting breaches for %s: %s", email, err)
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, breach.NewErrorf(breach.BreachNotFoundErr, "no breach found for %s", email)
	case http.StatusTooManyRequests:
		hrb := &hibp429ResponseBody{}
		defer resp.Body.Close()

		err := json.NewDecoder(resp.Body).Decode(hrb)
		if err != nil {
			return nil, breach.NewErrorf(breach.DatasourceErr, "could not parse 429 body: %v", err)
		}
		wt, err := getWaitTime(hrb)
		if err != nil {
			return nil, breach.NewErrorf(
				breach.DatasourceErr, "could not parse wait time from message: %v", err)
		}

		rresp, rerr := c.retry(req, email, wt)
		if rerr != nil {
			return nil, breach.NewErrorf(breach.DatasourceErr, "error in retrying: %s", rerr)
		}
		return rresp, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, breach.NewErrorf(
			breach.DatasourceErr, "error getting breaches for %s: %v", email, resp.StatusCode)
	}

	return resp, nil
}

// Sleep for 'waitTime' seconds before retrying.
func (c *hibpClient) retry(req *http.Request, email string, waitTime int) (*http.Response, *breach.Error) {
	if c.config.maxRetries == c.retries {
		return nil, breach.NewErrorf(
			breach.DatasourceErr, "too many retries: %v; max %v", c.retries, c.config.maxRetries)
	}
	wt := time.Duration(waitTime) * time.Second
	c.sleeper(wt)
	c.retries += 1

	return c.makeRequest(req, email)
}

func getWaitTime(b *hibp429ResponseBody) (int, error) {
	re, _ := regexp.Compile("[0-9]+")
	w8t := re.FindString(b.Message)
	if w8t == "" {
		return 0, fmt.Errorf("no number found in string: %v", b.Message)
	}
	iw8t, err := strconv.Atoi(w8t)
	if err != nil {
		return 0, fmt.Errorf("could not transform string number to int: %v", err)
	}
	return iw8t, nil
}

func (c *hibpClient) toBreach(hrb *hibp200ResponseBreach, email string) (*breach.Breach, error) {
	layout := "2006-01-02"
	bd, err := time.Parse(layout, hrb.BreachDate)
	if err != nil {
		return nil, err
	}
	b, err := breach.New(
		email,
		hrb.DataClasses,
		bd,
		hrb.Title,
	)
	if err != nil {
		return nil, err
	}

	return b, nil
}
