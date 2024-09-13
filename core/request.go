package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	apiURL     string
	referURL   string
	authToken  string
	httpClient *http.Client
}

func (c *Client) makeRequest(method string, endpoint string, jsonBody interface{}) ([]byte, error) {
	fullURL := c.apiURL + endpoint

	// Convert body to JSON
	var reqBody []byte
	var err error
	if jsonBody != nil {
		reqBody, err = json.Marshal(jsonBody)
		if err != nil {
			return nil, err
		}
	}

	// Create new request
	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// Set header
	setHeader(req, c.referURL, c.authToken)

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 status code
	if resp.StatusCode >= 400 {
		// Read the response body to include in the error message
		bodyBytes, bodyErr := io.ReadAll(resp.Body)
		if bodyErr != nil {
			return nil, fmt.Errorf("error status: %v, and failed to read body: %v", resp.StatusCode, bodyErr)
		}
		return nil, fmt.Errorf("error status: %v, error message: %s", resp.StatusCode, string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}

// Req Login
func (c *Client) loginAccount(account *Account) ([]byte, error) {
	// Struct for user object to maintain field order
	userPayload := struct {
		ID              int    `json:"id"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Username        string `json:"username"`
		LanguageCode    string `json:"language_code"`
		AllowsWriteToPm bool   `json:"allows_write_to_pm"`
	}{
		ID:              account.UserId,
		FirstName:       account.FirstName,
		LastName:        account.LastName,
		Username:        account.Username,
		LanguageCode:    account.LanguageCode,
		AllowsWriteToPm: account.AllowWriteToPm,
	}

	// Struct for the entire payload to maintain field order
	payload := struct {
		QueryID  string      `json:"query_id"`
		User     interface{} `json:"user"`
		AuthDate string      `json:"auth_date"`
		Hash     string      `json:"hash"`
	}{
		QueryID:  account.QueryId,
		User:     userPayload,
		AuthDate: account.AuthDate,
		Hash:     account.Hash,
	}

	return c.makeRequest("POST", "/login", payload)
}

// Req User Info
func (c *Client) userInfo() ([]byte, error) {
	return c.makeRequest("GET", "/auth/users/info", nil)
}

// Req Task List
func (c *Client) taskList() ([]byte, error) {
	return c.makeRequest("GET", "/auth/tasks/list", nil)
}

// Req Claim Task
func (c *Client) claimTask(taskName, taskType string) ([]byte, error) {
	payload := map[string]string{
		"name": taskName,
		"type": taskType,
	}

	return c.makeRequest("POST", "/auth/tasks/submit", payload)
}

// Req Achievement List
func (c *Client) achieveList() ([]byte, error) {
	return c.makeRequest("GET", "/auth/achievement/canClaim", nil)
}

// Req Claim Achievement
func (c *Client) claimAchievement(achieveName string) ([]byte, error) {
	return c.makeRequest("POST", fmt.Sprintf("/auth/achievement/claim/%s", achieveName), nil)
}

// Req Info Spin Wheel
func (c *Client) spinInfo() ([]byte, error) {
	return c.makeRequest("GET", "/auth/spins/current", nil)
}

// Req Spin Wheel
func (c *Client) spinWheel() ([]byte, error) {
	return c.makeRequest("POST", "/auth/spins/submit", nil)
}

// Req Info Staking
func (c *Client) stakingInfo() ([]byte, error) {
	return c.makeRequest("GET", "/auth/stakes/info", nil)
}

// Req Info Staking
func (c *Client) staking(amount string, poolId int) ([]byte, error) {
	payload := map[string]interface{}{
		"score":   amount,
		"pool_id": poolId,
	}

	return c.makeRequest("POST", "/auth/stakes/submit", payload)
}
