package cag

import (
	"testing"

	"net/http"

	"github.com/stretchr/testify/assert"
)

type fakeAlert struct{}

func (a *fakeAlert) Create() (AlertResponse, int, error) {
	return AlertResponse{}, http.StatusOK, nil
}

func (a *fakeAlert) Check() error {
	return nil
}

func TestNewAlert(t *testing.T) {
	assert := assert.New(t)
	client := &client{
		apiKey: fakeAPIKey,
	}
	alert := &AlertData{}
	client.NewAlert(alert)
	assert.Equal(client, alert.Client)
	assert.Equal("summary", alert.Updatable)
}

func TestCreateAlert(t *testing.T) {
	assert := assert.New(t)
	client := NewFakeClient()
	alert := AlertData{Client: client}
	alertResponse, statusCode, err := alert.Create()
	assert.NoError(err, "failed to create alert without an error")
	assert.Equal(fakeAlertResponseID, alertResponse.ID, "alert response does not match")
	assert.Equal(http.StatusCreated, statusCode, "http status not correct")
}
