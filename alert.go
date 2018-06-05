package cag

import (
	"encoding/json"
	"net/http"
)

// Return history alerts based on one or more alert paramaerts
//
// set Updatable to either:
// - summary (append to a single Spark INC for any given MonitoredItem + AlertSummary)
// - host (append to a single Spark INC for any given MonitoredItem)
// - false (new INC for each call)

// Alert interface
type Alert interface {
	// Create new cag alert
	Create() (AlertResponse, int, error)
	// Check alert is valid
	Check() error
}

// AlertData struct to send options
type AlertData struct {
	Client              Client
	MonitoredItem       string
	AlertSummary        string
	DetailedDescription string
	AssignTo            string
	Updatable           string
	AffectedArea        string
	AffectedCI          string
	Severity            string
	SparkInstance       string
	IncidentType        string
	HelpURL             string
	MonitoringGroup     string
	MonitoringSystem    string
}

// AlertResponse struct
type AlertResponse struct {
	ID string `json:"id"`
}

// NewAlert struct
func (c *client) NewAlert(a *AlertData) Alert {
	a.Client = c
	a.Updatable = "summary"
	return a
}

// Create method on alert
func (a *AlertData) Create() (AlertResponse, int, error) {
	var alertResponse AlertResponse
	if err := a.Check(); err != nil {
		return alertResponse, http.StatusBadRequest, err
	}
	kv := map[string]string{
		"monitored_item":       a.MonitoredItem,
		"alert_summary":        a.AlertSummary,
		"detailed_description": a.DetailedDescription,
		"assign_to":            a.AssignTo,
		"updatable":            a.Updatable,
		"affected_area":        a.AffectedArea,
		"affected_ci":          a.AffectedCI,
		"severity":             a.Severity,
		"spark_instance":       a.SparkInstance,
		"incident_type":        a.IncidentType,
		"help_url":             a.HelpURL,
		"monitoring_group":     a.MonitoringGroup,
		"monitoring_system":    a.MonitoringSystem,
	}
	responseCode, body, err := a.Client.Do("/v3/submit", "POST", kv)
	if err != nil {
		return alertResponse, responseCode, err
	}
	return alertResponse, responseCode, json.Unmarshal(body, &alertResponse)
}

// Check alert is valid
func (a *AlertData) Check() error {
	return nil
}
