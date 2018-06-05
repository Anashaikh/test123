package cag

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"encoding/json"

	"github.com/sethgrid/pester"
)

const cagFormat = "json"

// Client interface
type Client interface {
	// Status of cag endpoint
	Status() (int, error)
	// Heartbeat request to cag
	Heartbeat() (int, error)
	// Permissions request to cag
	Permissions() (PermissionsResponse, int, error)
	// NewAlert interface
	NewAlert(*AlertData) Alert
	// NewHistory interface
	NewHistory(*HistoryData) History
	// Do method to call cag
	Do(path, method string, extraURLValues map[string]string) (int, []byte, error)
}

// PermissionsResponse struct
type PermissionsResponse []struct {
	TeamName      string `json:"teamName"`
	SparkInstance string `json:"spark_instance"`
}

type client struct {
	serverURL  url.URL
	httpClient *pester.Client
	apiKey     string
}

// NewClient CAG client
func NewClient(serverURL string, apiKey string, sslSkipVerify bool) Client {
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
	c := &client{
		serverURL:  *parsedServerURL,
		httpClient: httpClient,
		apiKey:     apiKey,
	}
	return c
}

// Status method to check CAG health
func (c *client) Status() (int, error) {
	statusCode, _, err := c.Do("/v3", "GET", nil)
	return statusCode, err
}

// Heartbeat method to check CAG health
func (c *client) Heartbeat() (int, error) {
	statusCode, _, err := c.Do("/v3/heartbeat", "POST", nil)
	return statusCode, err
}

// Permissions method to check api key access
func (c *client) Permissions() (PermissionsResponse, int, error) {
	var permissionsResponse PermissionsResponse
	statusCode, body, err := c.Do("/v3/permissions", "GET", nil)
	if err != nil {
		return permissionsResponse, statusCode, err
	}
	return permissionsResponse, statusCode, json.Unmarshal(body, &permissionsResponse)
}

func (c *client) Do(path, method string, extraURLValues map[string]string) (int, []byte, error) {
	var body []byte
	u := c.serverURL
	u.Path = path
	kv := c.defaultURLValues()
	for key, value := range extraURLValues {
		addURLValue(&kv, key, value)
	}
	u.RawQuery = kv.Encode()
	req, _ := http.NewRequest(method, u.String(), nil)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, body, err
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return http.StatusInternalServerError, body, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusBadRequest {
		err = errors.New("required parameters missing: " + string(body))
	} else if res.StatusCode == http.StatusForbidden {
		err = errors.New("invalid api key")
	} else if res.StatusCode == http.StatusNotAcceptable {
		err = errors.New("alert already registered within 2 minutes")
	} else if res.StatusCode != http.StatusOK {
		err = errors.New("unrecognised cag status code: " + string(body))
	}
	return res.StatusCode, body, err
}

func (c *client) defaultURLValues() url.Values {
	urlValues := make(url.Values)
	addURLValue(&urlValues, "_format", cagFormat)
	addURLValue(&urlValues, "apikey", c.apiKey)
	return urlValues
}

func addURLValue(u *url.Values, k string, v string) {
	if len(v) > 0 {
		u.Add(k, v)
	}
}
