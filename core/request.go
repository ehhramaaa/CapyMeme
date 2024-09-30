package core

import (
	"fmt"
	"strings"
)

const apiUrl = "https://api.capybarameme.com"

// Req Login
func (c *Client) getToken() (string, error) {
	// Struct for user object to maintain field order
	userPayload := struct {
		ID              int    `json:"id"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Username        string `json:"username"`
		LanguageCode    string `json:"language_code"`
		AllowsWriteToPm bool   `json:"allows_write_to_pm"`
	}{
		ID:              c.account.userId,
		FirstName:       c.account.firstName,
		LastName:        c.account.lastName,
		Username:        c.account.username,
		LanguageCode:    c.account.languageCode,
		AllowsWriteToPm: c.account.allowWriteToPm,
	}

	// Struct for the entire payload to maintain field order
	payload := struct {
		QueryID  string      `json:"query_id"`
		User     interface{} `json:"user"`
		AuthDate string      `json:"auth_date"`
		Hash     string      `json:"hash"`
	}{
		QueryID:  c.account.queryId,
		User:     userPayload,
		AuthDate: c.account.authDate,
		Hash:     c.account.hash,
	}

	res, err := c.makeRequest("POST", "https://api.capybarameme.com/login", payload)
	if err != nil {
		return "", err
	}

	if token, exits := res["token"].(string); exits {
		return token, nil
	} else {
		return "", fmt.Errorf("Token not found!")
	}
}

// Req User Info
func (c *Client) getUserInfo() (map[string]interface{}, error) {
	res, err := c.makeRequest("GET", apiUrl+"/auth/users/info", nil)
	if err != nil {
		return nil, err
	}

	if user, exits := res["user"].(map[string]interface{}); exits {
		return user, nil
	} else {
		return nil, fmt.Errorf("Field user not found!")
	}
}

// Req Task List
func (c *Client) taskList() (map[string]interface{}, error) {
	res, err := c.makeRequest("GET", apiUrl+"/auth/tasks/list", nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Req Claim Task
func (c *Client) claimTask(taskName, taskType string) (string, error) {
	payload := map[string]string{
		"name": taskName,
		"type": taskType,
	}

	res, err := c.makeRequest("POST", apiUrl+"/auth/tasks/submit", payload)
	if err != nil {
		return "", err
	}
	if msg, exits := res["msg"].(string); exits && strings.Contains(msg, "completed") {
		return msg, nil
	} else {
		return "", fmt.Errorf("Task not completed!")
	}
}

// Req Achievement List
func (c *Client) achievementList() (map[string]interface{}, error) {
	res, err := c.makeRequest("GET", apiUrl+"/auth/achievement/canClaim", nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Req Claim Achievement
func (c *Client) claimAchievement(achieveName string) (string, error) {
	res, err := c.makeRequest("POST", fmt.Sprintf("/auth/achievement/claim/%s", achieveName), nil)
	if err != nil {
		return "", err
	}
	if msg, exits := res["response"].(string); exits {
		return msg, nil
	} else {
		return "", fmt.Errorf("Achievement not completed!")
	}
}

// Req Info Spin Wheel
func (c *Client) spinInfo() (map[string]interface{}, error) {
	res, err := c.makeRequest("GET", apiUrl+"/auth/spins/current", nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Req Spin Wheel
func (c *Client) spinWheel() (map[string]interface{}, error) {
	res, err := c.makeRequest("POST", apiUrl+"/auth/spins/submit", nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Req Info Staking
func (c *Client) stakingInfo() (map[string]interface{}, error) {
	res, err := c.makeRequest("GET", apiUrl+"/auth/stakes/info", nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Req Info Staking
func (c *Client) staking(amount string, poolId int) {
	payload := map[string]interface{}{
		"pool_id": poolId,
		"score":   amount,
	}

	c.makeRequest("POST", apiUrl+"/auth/stakes/submit", payload)
}
