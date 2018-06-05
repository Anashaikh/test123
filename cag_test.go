package cag

import (
	"net/http"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

const (
	fakeAPIKey          = "abcd1234"
	fakeURL             = "https://test-url"
	fakeSSLSkipVerify   = false
	fakeAlertResponseID = "99"
)

type fakeClient struct{}

func NewFakeClient() Client {
	return &fakeClient{}
}

func (c *fakeClient) Status() (int, error) {
	return http.StatusOK, nil
}

func (c *fakeClient) Heartbeat() (int, error) {
	return http.StatusOK, nil
}

func (c *fakeClient) Permissions() (PermissionsResponse, int, error) {
	return PermissionsResponse{}, http.StatusOK, nil
}

func (c *fakeClient) NewAlert(*AlertData) Alert {
	return &fakeAlert{}
}

func (c *fakeClient) NewHistory(*HistoryData) History {
	return &fakeHistory{}
}

func (c *fakeClient) Do(path, method string, extraURLValues map[string]string) (int, []byte, error) {
	var statusCode int
	var body []byte
	var err error
	switch path {
	case "/v3/submit":
		statusCode = http.StatusCreated
		body = []byte(fmt.Sprintf("{\"id\":\"%s\"}", fakeAlertResponseID))
	case "/v3/history":
		statusCode = http.StatusFound
		body = []byte(`{"total_pages":5}`)
	}
	return statusCode, body, err
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)
	client := NewClient(fakeURL, fakeAPIKey, fakeSSLSkipVerify)
	assert.Implements(new(Client), client)
}

func TestDefaultURLValues(t *testing.T) {
	assert := assert.New(t)
	client := &client{apiKey: fakeAPIKey}
	urlValues := client.defaultURLValues()
	assert.Equal("_format="+cagFormat+"&apikey="+fakeAPIKey, urlValues.Encode(),
		"default url values does not create expected value when api key is present")
}

func TestDefaultURLValuesEmptyAPIKey(t *testing.T) {
	assert := assert.New(t)
	client := &client{}
	urlValues := client.defaultURLValues()
	assert.Equal("_format="+cagFormat, urlValues.Encode(),
		"default url values does not create expected value when api key is omitted")
}
