package core

import (
	"CapybaraMeme/tools"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/gookit/config/v2"
)

func (account *Account) parsingQueryData() {
	value, err := url.ParseQuery(account.queryData)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed to parse query data: %s", err))
	}

	if len(value.Get("query_id")) > 0 {
		account.queryId = value.Get("query_id")
	}

	if len(value.Get("auth_date")) > 0 {
		account.authDate = value.Get("auth_date")
	}

	if len(value.Get("hash")) > 0 {
		account.hash = value.Get("hash")
	}

	userParam := value.Get("user")

	var userData map[string]interface{}
	err = json.Unmarshal([]byte(userParam), &userData)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed to parse user data: %s", err))
	}

	userId, ok := userData["id"].(float64)
	if !ok {
		tools.Logger("error", "Failed to convert ID to float64")
	}

	account.userId = int(userId)

	username, ok := userData["username"].(string)
	if !ok {
		tools.Logger("error", "Failed to get username from query")
		return
	}

	account.username = username

	// Ambil first name
	firstName, ok := userData["first_name"].(string)
	if !ok {
		tools.Logger("error", "Failed to get first name from query")
	}

	account.firstName = firstName

	// Ambil first name
	lastName, ok := userData["last_name"].(string)
	if !ok {
		tools.Logger("error", "Failed to get last name from query")
	}
	account.lastName = lastName

	// Ambil language code
	languageCode, ok := userData["language_code"].(string)
	if !ok {
		tools.Logger("error", "Failed to get language code from query")
	}
	account.languageCode = languageCode

	// Ambil allowWriteToPm
	allowWriteToPm, ok := userData["allows_write_to_pm"].(bool)
	if !ok {
		tools.Logger("error", "Failed to get allows write to pm from query")
	}

	account.allowWriteToPm = allowWriteToPm
}

func (c *Client) autoCompleteTask() int {
	var points int
	isAutoSpin := config.Bool("AUTO_SPIN")
	isAutoStake := config.Bool("AUTO_STAKING")
	targetStaking := config.Int("STAKING_PERCENTAGE")

	token, err := c.getToken()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get token: %v", c.account.username, err))
	}

	if len(token) > 0 {
		c.accessToken = fmt.Sprintf("Bearer %s", token)
	} else {
		tools.Logger("error", fmt.Sprintf("| %s | Token not found!", c.account.username))
		return points
	}

	userInfo, err := c.getUserInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get user info: %v", c.account.username, err))
		return points
	}

	if userInfo != nil {
		tools.Logger("success", fmt.Sprintf("| %s | Total Balance: %s | Stacked: %s", c.account.username, userInfo["total_score"].(string), userInfo["locked_score"].(string)))
	}

	taskLists, err := c.taskList()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get task list: %v", c.account.username, err))
	}

	for _, list := range taskLists {
		if listMap, exits := list.(map[string]interface{}); exits {
			if tasks, exists := listMap["list"].([]interface{}); exists {
				for _, task := range tasks {
					if taskMap, exits := task.(map[string]interface{}); exits {
						if !taskMap["is_completed"].(bool) {
							claimTask, err := c.claimTask(taskMap["name"].(string), listMap["type"].(string))
							if err != nil {
								tools.Logger("error", fmt.Sprintf("| %s | Failed To Claim Task : %v | Sleep 5s Before Claim Another Task...", c.account.username, err))
							}

							if claimTask != "" {
								tools.Logger("success", fmt.Sprintf("| %s | Claim Task %s Successfully | Sleep 5s Before Claim Another Task...", c.account.username, taskMap["name"].(string)))

							}

							time.Sleep(5 * time.Second)
						}
					}
				}
			}
		}
	}

	achievementList, err := c.achievementList()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get achievement list: %v", c.account.username, err))
	}

	if achievementList != nil {
		for _, achieve := range achievementList {
			if achieveMap, exits := achieve.(map[string]interface{}); exits {
				if list, exits := achieveMap["achievements"].([]interface{}); exits && len(list) > 0 {
					for _, achievement := range list {
						if achievementInfo := achievement.(map[string]interface{}); exits {
							claimAchievement, err := c.claimAchievement(achievementInfo["name"].(string))
							if err != nil {
								tools.Logger("error", fmt.Sprintf("| %s | Failed To Claim Achievement %s : %v | Sleep 5s Before Claim Another Achievement...", c.account.username, achievementInfo["name"].(string), err))
							}

							if claimAchievement != "" {
								tools.Logger("success", fmt.Sprintf("| %s | Claim Achievement %s %s | Sleep 5s Before Claim Another Achievement...", c.account.username, achieve, claimAchievement))
							}

							time.Sleep(5 * time.Second)
						}
					}
				}
			}
		}
	}

	if isAutoSpin {
		spinInfo, err := c.spinInfo()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to get spin info: %v", c.account.username, err))
		}

		if spinInfo != nil {
			availableSpin := int((spinInfo["free_spins"].(float64) + spinInfo["remain_spins"].(float64)))

			if rewardBalance, exits := spinInfo["acc_rewards"].(map[string]interface{}); exits {
				tools.Logger("success", fmt.Sprintf("| %s | Total Balance: %s | Stacked: %s | Ton: %s | Gem: %s | Spin: %v", c.account.username, userInfo["total_score"].(string), userInfo["locked_score"].(string), rewardBalance["ton"].(string), rewardBalance["gem"].(string), availableSpin))
			}

			if availableSpin > 0 {
				for availableSpin > 0 {
					spinWheel, err := c.spinWheel()
					if err != nil {
						tools.Logger("error", fmt.Sprintf("| %s | Failed to spin wheel: %v | Sleep 5s Before Next Spin...", c.account.username, err))
					}

					if spinWheel != nil {
						if reward, exits := spinWheel["reward"].(map[string]interface{}); exits {
							if spinInfo, exits = spinWheel["task"].(map[string]interface{}); exits {
								availableSpin = int((spinInfo["free_spins"].(float64) + spinInfo["remain_spins"].(float64)))

								tools.Logger("success", fmt.Sprintf("| %s | Spin Wheel Successfully | Reward: %s %s | Spin Balance: %v | Sleep 5s Before Next Spin...", c.account.username, reward["reward_value"].(string), reward["reward_type"].(string), availableSpin))
							}
						}
					}

					time.Sleep(5 * time.Second)
				}
			}
		}
	}

	if isAutoStake {
		totalPoints, _ := strconv.Atoi(userInfo["total_score"].(string))
		stacked, _ := strconv.Atoi(userInfo["locked_score"].(string))
		stakingPercentage := (float64(stacked) / float64(totalPoints)) * 100

		if stakingPercentage < float64(targetStaking) {
			stakingInfo, err := c.stakingInfo()
			if err != nil {
				tools.Logger("error", fmt.Sprintf("| %s | Failed to get staking info: %v", c.account.username, err))
			}

			if stakingInfo != nil {

				stakingPool := stakingInfo["pool"].(map[string]interface{})

				remainingPercentage := float64(targetStaking) - stakingPercentage

				amountToStake := fmt.Sprintf("%.f", (remainingPercentage/100)*float64(totalPoints))

				c.staking(amountToStake, int(stakingPool["pool_id"].(float64)))

				updateStakingInfo, err := c.stakingInfo()
				if err != nil {
					tools.Logger("error", fmt.Sprintf("| %s | Failed to get update staking info: %v", c.account.username))
				}

				if stake, exits := stakingInfo["stake"].(map[string]interface{}); exits {
					if update, exits := updateStakingInfo["stake"].(map[string]interface{}); exits {
						if update["stake_score"].(string) != stake["stake_score"].(string) {
							tools.Logger("success", fmt.Sprintf("| %s | Staking %s Score Successfully  | Current Balance: %.0f | Total Stacked: %s | Realized: %.0f | Pending: %.0f | Total Reward: %.0f | Claimable: %v", c.account.username, amountToStake, update["available_score"].(float64), update["stake_score"].(string), update["realized"].(float64), update["pending"].(float64), update["total_reward"].(float64), update["has_claimables"].(bool)))
						} else {
							tools.Logger("error", fmt.Sprintf("| %s | Failed to staking %s Score...", c.account.username, amountToStake))
						}
					}
				}
			}
		}
	}

	userInfo, err = c.getUserInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get user info: %v", c.account.username, err))
		return points
	}

	if userInfo != nil {
		points, _ = strconv.Atoi(userInfo["total_score"].(string))
	}

	return points
}
