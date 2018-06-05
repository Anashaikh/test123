package cag

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Return history alerts based on one or more alert paramaerts

// History struct to send options
type History struct {
	Client              *Client
	MonitoredItem       string
	AlertSummary        string
	DetailedDescription string
	AffectedArea        string
	AffectedCI          string
	Severity            string
	SparkInstance       string
	IncidentType        string
	HelpURL             string
	MonitoringGroup     string
	MonitoringSystem    string
	Page                int
}

// HistoryResponse struct
type HistoryResponse struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	Alerts      []struct {
		ID                  int    `json:"id"`
		SparkInstance       string `json:"spark_instance"`
		IncidentType        string `json:"incident_type"`
		MonitoredItem       string `json:"monitoredItem"`
		AlertSummary        string `json:"alertSummary"`
		DetailedDescription string `json:"detailedDescription"`
		HelpURL             string `json:"help_url"`
		AssignTo            string `json:"assignTo"`
		Updatable           string `json:"updatable"`
		AffectedArea        string `json:"affectedArea"`
		AffectedCi          string `json:"affectedCi"`
		Severity            string `json:"severity"`
	} `json:"alerts"`
}

// NewHistory struct
func (c *Client) NewHistory() *History {
	return &History{
		Client: c,
		Page:   1,
	}
}

// CheckHistory checks histroy stuct for valid inputs
func (h *History) CheckHistory() error {
	if h.Page <= 0 {
		return errors.New("Page Number Must Be Non-Negative, got value: " + strconv.Itoa(h.Page))
	}
	return nil
}

// Get method for History
func (h *History) Get() (HistoryResponse, int, error) {
	var historyResponse HistoryResponse
	if err := h.CheckHistory(); err != nil {
		return historyResponse, http.StatusBadRequest, err
	}
	url := h.Client.ServerURL
	url.Path = "/v3/history"
	urlValues := h.Client.DefaultURLValues()
	addURLValue(&urlValues, "monitored_item", h.MonitoredItem)
	addURLValue(&urlValues, "alert_summary", h.AlertSummary)
	addURLValue(&urlValues, "detailed_description", h.DetailedDescription)
	addURLValue(&urlValues, "affected_area", h.AffectedArea)
	addURLValue(&urlValues, "affected_ci", h.AffectedCI)
	addURLValue(&urlValues, "severity", h.Severity)
	addURLValue(&urlValues, "spark_instance", h.SparkInstance)
	addURLValue(&urlValues, "incident_type", h.IncidentType)
	addURLValue(&urlValues, "help_url", h.HelpURL)
	addURLValue(&urlValues, "monitoring_group", h.MonitoringGroup)
	addURLValue(&urlValues, "monitoring_system", h.MonitoringSystem)
	addURLValue(&urlValues, "page", strconv.Itoa(h.Page))
	url.RawQuery = urlValues.Encode()
	req, _ := http.NewRequest("GET", url.String(), nil)
	res, err := h.Client.HTTPClient.Do(req)
	if err != nil {
		return historyResponse, http.StatusInternalServerError, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return historyResponse, http.StatusInternalServerError, err
	}
	defer res.Body.Close()
	if res.StatusCode == 204 {
		// No alerts submitted for given API Key, return empty struct
		return historyResponse, res.StatusCode, nil
	} else if res.StatusCode == 403 {
		return historyResponse, res.StatusCode, errors.New("Invalid API Key")
	} else if res.StatusCode != 200 {
		return historyResponse, res.StatusCode, errors.New("Unrecognised CAG Status Code: " + string(body))
	}
	return historyResponse, res.StatusCode, json.Unmarshal(body, &historyResponse)
}
