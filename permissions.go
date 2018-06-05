package cag

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Permissions struct
type Permissions []struct {
	TeamName      string `json:"teamName"`
	SparkInstance string `json:"spark_instance"`
}

// Permissions method to check CAG health
func (c *Client) Permissions() (Permissions, int, error) {
	var permissions Permissions
	url := c.ServerURL
	url.Path = "/v3/permissions"
	url.RawQuery = c.DefaultURLValues().Encode()
	req, _ := http.NewRequest("GET", url.String(), nil)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return permissions, http.StatusInternalServerError, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return permissions, http.StatusInternalServerError, err
	}
	defer res.Body.Close()
	if res.StatusCode == 403 {
		return permissions, res.StatusCode, errors.New("Invalid API Key")
	} else if res.StatusCode != 200 {
		return permissions, res.StatusCode, errors.New("Unrecognised CAG Status Code: " + string(body))
	}
	return permissions, res.StatusCode, json.Unmarshal(body, &permissions)
}
