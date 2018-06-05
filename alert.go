package cag

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Return history alerts based on one or more alert paramaerts
//
// set Updatable to either:
// - summary (append to a single Spark INC for any given MonitoredItem + AlertSummary)
// - host (append to a single Spark INC for any given MonitoredItem)
// - false (new INC for each call)

// Alert struct to send options
type Alert struct {
	Client              *Client
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
func (c *Client) NewAlert() *Alert {
	return &Alert{
		Client:    c,
		Updatable: "summary",
	}
}

// CheckAlert checks alert struct to have valid inputs
func (a *Alert) CheckAlert() error {
	return nil
}

// Create method on alert
func (a *Alert) Create() (AlertResponse, int, error) {
	var alertResponse AlertResponse
	if err := a.CheckAlert(); err != nil {
		return alertResponse, http.StatusBadRequest, err
	}
	url := a.Client.ServerURL
	url.Path = "/v3/submit"
	urlValues := a.Client.DefaultURLValues()
	addURLValue(&urlValues, "monitored_item", a.MonitoredItem)
	addURLValue(&urlValues, "alert_summary", a.AlertSummary)
	addURLValue(&urlValues, "detailed_description", a.DetailedDescription)
	addURLValue(&urlValues, "assign_to", a.AssignTo)
	addURLValue(&urlValues, "updatable", a.Updatable)
	addURLValue(&urlValues, "affected_area", a.AffectedArea)
	addURLValue(&urlValues, "affected_ci", a.AffectedCI)
	addURLValue(&urlValues, "severity", a.Severity)
	addURLValue(&urlValues, "spark_instance", a.SparkInstance)
	addURLValue(&urlValues, "incident_type", a.IncidentType)
	addURLValue(&urlValues, "help_url", a.HelpURL)
	addURLValue(&urlValues, "monitoring_group", a.MonitoringGroup)
	addURLValue(&urlValues, "monitoring_system", a.MonitoringSystem)
	url.RawQuery = urlValues.Encode()
	req, _ := http.NewRequest("POST", url.String(), nil)
	res, err := a.Client.HTTPClient.Do(req)
	if err != nil {
		return alertResponse, http.StatusInternalServerError, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return alertResponse, http.StatusInternalServerError, err
	}
	defer res.Body.Close()
	if res.StatusCode == 400 {
		return alertResponse, res.StatusCode, errors.New("Required Parameters Missing: " + string(body))
	} else if res.StatusCode == 403 {
		return alertResponse, res.StatusCode, errors.New("Invalid API Key")
	} else if res.StatusCode == 406 {
		return alertResponse, res.StatusCode, errors.New("Alert Already Registered Within 2 Minutes")
	} else if res.StatusCode != 200 {
		return alertResponse, res.StatusCode, errors.New("Unrecognised CAG Status Code: " + string(body))
	}
	return alertResponse, res.StatusCode, json.Unmarshal(body, &alertResponse)
}
