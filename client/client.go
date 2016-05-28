// Package client is an http client for the userauth service
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeffbmartinez/userauth/model"
)

const (
	sessionInfoEndpointTemplate = "%s://%s:%s/session/info"
	pingEndpointTemplate        = "%s://%s:%s/ping"
)

// Client is a userauth service client instance.
type Client struct {
	Protocol string
	Hostname string
	Port     string
}

// NewClient creates an instance of the userauth client.
func NewClient(hostname string, port string) *Client {
	return &Client{
		Protocol: "http",
		Hostname: hostname,
		Port:     port,
	}
}

// SessionInfoRequestBody represents the body of a request to /session/info
type SessionInfoRequestBody struct {
	SessionInfo string `json:"sessionInfo"`
}

// SessionInfo calls the user auth endpoint to decode an encrypted session ID cookie
// and returns the contents. If err is set, something went wrong. If err is nil but
// the cookie returned is empty, the session is not valid.
func (c Client) SessionInfo(sessionInfo string) (model.SessionCookie, error) {
	sessionInfoRequestBody := SessionInfoRequestBody{
		SessionInfo: sessionInfo,
	}

	sessionInfoRequestBodyBytes, err := json.Marshal(sessionInfoRequestBody)
	if err != nil {
		return model.SessionCookie{}, err
	}

	requestBody := bytes.NewBuffer(sessionInfoRequestBodyBytes)
	endpoint := fmt.Sprintf(sessionInfoEndpointTemplate, c.Protocol, c.Hostname, c.Port)
	response, err := http.Post(endpoint, "application/json", requestBody)
	if err != nil {
		return model.SessionCookie{}, err
	}

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("Response code not OK, got %d (%s)", response.StatusCode, response.Status)
		return model.SessionCookie{}, err
	}

	var sessionCookie model.SessionCookie
	err = json.NewDecoder(response.Body).Decode(&sessionCookie)
	if err != nil {
		return model.SessionCookie{}, err
	}
	defer response.Body.Close()

	return sessionCookie, nil
}

// Ping returns an error if there was a problem pinging the user auth service.
// On successful ping, Ping returns nil.
func (c Client) Ping() error {
	endpoint := fmt.Sprintf(pingEndpointTemplate, c.Protocol, c.Hostname, c.Port)
	response, err := http.Get(endpoint)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected 200 OK response from /ping, got %d %s", response.StatusCode, response.Status)
	}

	return nil
}
