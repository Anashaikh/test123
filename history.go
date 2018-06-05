package cag

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// Return history alerts based on one or more alert paramaerts

// History interface
type History interface {
	// Get history from cag
	Get() (HistoryResponse, int, error)
	// Check history request is valid
	Check() error
}

// HistoryData struct to send options
type HistoryData struct {
	Client              Client
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
func (c *client) NewHistory(h *HistoryData) History {
	h.Client = c
	h.Page = 1
	return h
}

// Get method for History
func (h *HistoryData) Get() (HistoryResponse, int, error) {
	var historyResponse HistoryResponse
	if err := h.Check(); err != nil {
		return historyResponse, http.StatusBadRequest, err
	}
	kv := map[string]string{
		"monitored_item":       h.MonitoredItem,
		"alert_summary":        h.AlertSummary,
		"detailed_description": h.DetailedDescription,
		"affected_area":        h.AffectedArea,
		"affected_ci":          h.AffectedCI,
		"severity":             h.Severity,
		"spark_instance":       h.SparkInstance,
		"incident_type":        h.IncidentType,
		"help_url":             h.HelpURL,
		"monitoring_group":     h.MonitoringGroup,
		"monitoring_system":    h.MonitoringSystem,
		"page":                 strconv.Itoa(h.Page),
	}
	statusCode, body, err := h.Client.Do("/v3/history", "GET", kv)
	if err != nil {
		return historyResponse, statusCode, err
	}
	return historyResponse, statusCode, json.Unmarshal(body, &historyResponse)
}

// Check history is valid
func (h *HistoryData) Check() error {
	if h.Page <= 0 {
		return errors.New("page number must be non-negative, got value: " + strconv.Itoa(h.Page))
	}
	return nil
}
