package cag

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeHistory struct{}

func (a *fakeHistory) Get() (HistoryResponse, int, error) {
	return HistoryResponse{}, http.StatusOK, nil
}

func (a *fakeHistory) Check() error {
	return nil
}

func TestNewHistory(t *testing.T) {
	assert := assert.New(t)
	client := &client{
		apiKey: fakeAPIKey,
	}
	history := &HistoryData{}
	client.NewHistory(history)
	assert.Equal(client, history.Client)
	assert.Equal(1, history.Page)
}

func TestGetHistory(t *testing.T) {
	assert := assert.New(t)
	client := NewFakeClient()
	history := HistoryData{Client: client, Page: 1}
	historyResponse, statusCode, err := history.Get()
	assert.NoError(err, "failed to get history without an error")
	assert.Equal(5, historyResponse.TotalPages, "alert response does not match")
	assert.Equal(http.StatusFound, statusCode, "http status not correct")
}

func TestGetHistoryInvalid(t *testing.T) {
	assert := assert.New(t)
	client := NewFakeClient()
	history := HistoryData{Client: client}
	_, statusCode, err := history.Get()
	assert.Error(err, "failed to get an error with invalid input")
	assert.Equal(http.StatusBadRequest, statusCode, "http status not correct")
}
