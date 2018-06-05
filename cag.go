package cag

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/sethgrid/pester"
)

const cagFormat = "json"

// Client struct
type Client struct {
	ServerURL  url.URL
	HTTPClient *pester.Client
	APIKey     string
}

// NewClient CAG client
func NewClient(serverURL string, apiKey string, sslSkipVerify bool) *Client {
	// Create pester client with default settings
	httpClient := pester.New()
	httpClient.Transport = &http.Transport{
		IdleConnTimeout: 15 * time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: sslSkipVerify},
	}
	httpClient.MaxRetries = 5
	httpClient.Backoff = pester.ExponentialBackoff
	httpClient.KeepLog = true
	httpClient.Timeout = 5 * time.Second
	// Parse serverURL and setup default url values
	parsedServerURL, _ := url.Parse(serverURL)
	c := &Client{
		ServerURL:  *parsedServerURL,
		HTTPClient: httpClient,
		APIKey:     apiKey,
	}
	return c
}

// Status method to check CAG health
func (c *Client) Status() (int, error) {
	url := c.ServerURL
	url.Path = "/v3"
	url.RawQuery = c.DefaultURLValues().Encode()
	req, _ := http.NewRequest("GET", url.String(), nil)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return res.StatusCode, errors.New("Unrecognised CAG Status Code: " + string(body))
	}
	return res.StatusCode, nil
}

// Heartbeat method to check CAG health
func (c *Client) Heartbeat() (int, error) {
	url := c.ServerURL
	url.Path = "/v3/heartbeat"
	url.RawQuery = c.DefaultURLValues().Encode()
	req, _ := http.NewRequest("POST", url.String(), nil)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer res.Body.Close()
	if res.StatusCode == 403 {
		return res.StatusCode, errors.New("Invalid API Key")
	} else if res.StatusCode != 200 {
		return res.StatusCode, errors.New("Unrecognised CAG Status Code: " + string(body))
	}
	return res.StatusCode, nil
}

// DefaultURLValues get default URL Values
func (c *Client) DefaultURLValues() url.Values {
	urlValues := make(url.Values)
	addURLValue(&urlValues, "_format", "json")
	addURLValue(&urlValues, "apikey", c.APIKey)
	return urlValues
}

func addURLValue(u *url.Values, k string, v string) {
	if len(v) > 0 {
		u.Add(k, v)
	}
}
