package core

import (
	"fmt"
	"time"

	"capybara-meme/helper"
)

func (c *Client) getToken(account *Account) map[string]interface{} {
	req, err := c.loginAccount(account)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to login: %v", err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Error handling response: %v", err))
		return nil
	}

	return res
}

func (c *Client) getUserInfo() map[string]interface{} {
	req, err := c.userInfo()
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to get user info: %v", err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Error handling response: %v", err))
		return nil
	}

	return res
}

func (c *Client) getTaskList() []map[string]interface{} {
	req, err := c.taskList()
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to get task list: %v", err))
		return nil
	}

	res, err := handleResponseArray(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Error handling response: %v", err))
		return nil
	}

	return res
}

func (c *Client) getAchieveList() map[string]interface{} {
	req, err := c.achieveList()
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to get achievement list: %v", err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Error handling response: %v", err))
		return nil
	}

	return res
}

func (c *Client) getSpinInfo() map[string]interface{} {
	req, err := c.spinInfo()
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to get spin info: %v", err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Error handling response: %v", err))
		return nil
	}

	return res
}

func (c *Client) getStakingInfo() map[string]interface{} {
	req, err := c.stakingInfo()
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to get staking info: %v", err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Error handling response: %v", err))
		return nil
	}

	return res
}

func (c *Client) completingTask(taskName, taskType string) map[string]interface{} {
	req, err := c.claimTask(taskName, taskType)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to completing task %s: %v", taskName, err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Error handling response: %v", err))
		return nil
	}

	return res
}

func (c *Client) autoSpinWheel() map[string]interface{} {
	req, err := c.spinWheel()
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to spin wheel: %v", err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Error handling response: %v", err))
		return nil
	}

	return res
}

func (c *Client) autoStaking(amount string, poolId int) map[string]interface{} {
	_, err := c.staking(amount, poolId)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to staking: %v", err))
		return nil
	}

	time.Sleep(5 * time.Second)

	stakingInfo := c.getStakingInfo()

	if info, exits := stakingInfo["stake"].(map[string]interface{}); exits {
		return info
	}

	helper.PrettyLog("error", "Failed to get info staking after auto staking")

	return nil
}
