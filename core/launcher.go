package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gookit/config/v2"

	"capybara-meme/helper"
)

type Account struct {
	QueryId        string
	UserId         int
	Username       string
	FirstName      string
	LastName       string
	AuthDate       string
	Hash           string
	AllowWriteToPm bool
	LanguageCode   string
	QueryData      string
}

func getAccountFromQuery(account *Account) {
	// Parsing Query To Get Username
	value, err := url.ParseQuery(account.QueryData)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to parse query: %v", err.Error()))
		return
	}

	if len(value.Get("query_id")) > 0 {
		account.QueryId = value.Get("query_id")
	}

	if len(value.Get("auth_date")) > 0 {
		account.AuthDate = value.Get("auth_date")
	}

	if len(value.Get("hash")) > 0 {
		account.Hash = value.Get("hash")
	}

	userParam := value.Get("user")

	// Mendekode string JSON
	var userData map[string]interface{}
	err = json.Unmarshal([]byte(userParam), &userData)
	if err != nil {
		panic(err)
	}

	// Mengambil ID dan username dari hasil decode
	userIDFloat, ok := userData["id"].(float64)
	if !ok {
		helper.PrettyLog("error", "Failed to convert ID to float64")
		return
	}

	account.UserId = int(userIDFloat)

	// Ambil username
	username, ok := userData["username"].(string)
	if !ok {
		helper.PrettyLog("error", "Failed to get username")
		return
	}
	account.Username = username

	// Ambil first name
	firstName, ok := userData["first_name"].(string)
	if !ok {
		helper.PrettyLog("error", "Failed to get first_name")
		return
	}
	account.FirstName = firstName

	// Ambil first name
	lastName, ok := userData["last_name"].(string)
	if !ok {
		helper.PrettyLog("error", "Failed to get last_name")
		return
	}
	account.LastName = lastName

	// Ambil language code
	languageCode, ok := userData["language_code"].(string)
	if !ok {
		helper.PrettyLog("error", "Failed to get language_code")
		return
	}
	account.LanguageCode = languageCode

	// Ambil allowWriteToPm
	allowWriteToPm, ok := userData["allows_write_to_pm"].(bool)
	if !ok {
		helper.PrettyLog("error", "Failed to get allows_write_to_pm")
		return
	}
	account.AllowWriteToPm = allowWriteToPm
}

func ProcessBot(config *config.Config) {
	queryPath := config.String("query-file")
	apiUrl := config.String("bot.api-url")
	referUrl := config.String("bot.refer-url")
	maxThread := config.Int("max-thread")
	isAutoSpin := config.Bool("auto-spin")
	isAutoStake := config.Bool("auto-stake")
	amountStake := config.Float("amount-stake")

	if int(amountStake) > 1 {
		helper.PrettyLog("error", "Amount Stake Must 0 - 1, Example : 1 = 100%, 0.01 = 1%, 0 = 0%")
		return
	}

	queryData := helper.ReadFileTxt(queryPath)
	if queryData == nil {
		helper.PrettyLog("error", "Query data not found")
		return
	}

	helper.PrettyLog("info", fmt.Sprintf("%v Query Data Detected", len(queryData)))
	helper.PrettyLog("info", "Start Processing Account...")

	time.Sleep(3 * time.Second)

	var wg sync.WaitGroup

	// Membuat semaphore dengan buffered channel
	semaphore := make(chan struct{}, maxThread)

	for j, query := range queryData {
		wg.Add(1)

		// Goroutine untuk setiap job
		go func(index int, query string) {
			defer wg.Done()

			// Mengambil token dari semaphore sebelum menjalankan job
			semaphore <- struct{}{}

			client := &Client{
				apiURL:   apiUrl,
				referURL: referUrl,
			}

			account := &Account{
				QueryData: query,
			}

			getAccountFromQuery(account)

			helper.PrettyLog("info", fmt.Sprintf("%s | Started Bot...", account.Username))

			// Jalankan bot
			launchBot(client, account, isAutoSpin, isAutoStake, amountStake)

			// Sleep setelah job selesai
			randomSleep := helper.RandomNumber(config.Int("random-sleep.min"), config.Int("random-sleep.max"))

			helper.PrettyLog("info", fmt.Sprintf("%s | Launch Bot Finished, Sleeping for %v seconds..", account.Username, randomSleep))

			// Melepaskan token dari semaphore
			<-semaphore

			time.Sleep(time.Duration(randomSleep) * time.Second)
		}(j, query)
	}

	// Tunggu sampai semua worker selesai memproses pekerjaan
	wg.Wait()

	// Program utama berjalan terus menerus
	select {} // Block forever to keep the program running
}

func launchBot(client *Client, account *Account, isAutoSpin bool, isAutoStake bool, amountStake float64) {
	client.httpClient = &http.Client{}

	userData := client.getToken(account)

	if len(userData["token"].(string)) > 0 {
		client.authToken = userData["token"].(string)
	} else {
		return
	}

	userData = client.getUserInfo()

	if user, exits := userData["user"].(map[string]interface{}); exits {
		spinInfo := client.getSpinInfo()

		if rewardBalance, exits := spinInfo["acc_rewards"].(map[string]interface{}); exits {
			helper.PrettyLog("success", fmt.Sprintf("%s | Total Balance: %s | Stacked: %s | Ton: %s | Gem: %s | Spin: %.0f", account.Username, user["total_score"].(string), user["locked_score"].(string), rewardBalance["ton"].(string), rewardBalance["gem"].(string), (spinInfo["free_spins"].(float64)+spinInfo["remain_spins"].(float64))))
		}
	} else {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed To Get User Data...", account.Username))
		return
	}

	taskLists := client.getTaskList()

	for _, list := range taskLists {
		if tasks, exists := list["list"].([]interface{}); exists {
			for _, task := range tasks {
				if taskMap, ok := task.(map[string]interface{}); ok {
					if !taskMap["is_completed"].(bool) {
						result := client.completingTask(taskMap["name"].(string), list["type"].(string))

						if res, exits := result["msg"].(string); exits && strings.Contains(res, "completed") {
							helper.PrettyLog("success", fmt.Sprintf("%s | Claim Task %s Successfully | Sleep 5s Before Claim Another Task...", account.Username, taskMap["name"].(string)))

							time.Sleep(5 * time.Second)
						} else {
							helper.PrettyLog("error", fmt.Sprintf("%s | Claim Task %s Failed | Sleep 5s Before Claim Another Task...", account.Username, taskMap["name"].(string)))
							time.Sleep(5 * time.Second)
						}
					}
				}
			}
		}
	}

	achieveList := client.getAchieveList()

	if list, exits := achieveList["achievements"].([]string); exits && len(list) > 0 {
		for _, achieve := range list {
			result, err := client.claimAchievement(achieve)
			if err != nil {
				helper.PrettyLog("error", fmt.Sprintf("Failed To Claim Achievement : %v", err))
				continue
			}

			helper.PrettyLog("success", fmt.Sprintf("%s | Claim Achievement %s %s | Sleep 5s Before Claim Another Achievement...", account.Username, achieve, string(result)))

			time.Sleep(5 * time.Second)
		}
	}

	if isAutoSpin {
		spinInfo := client.getSpinInfo()

		limit := false

		if int((spinInfo["free_spins"].(float64) + spinInfo["remain_spins"].(float64))) > 0 {
			for !limit {
				spinWheel := client.autoSpinWheel()

				if reward, exits := spinWheel["reward"].(map[string]interface{}); exits {
					spinInfo = spinWheel["task"].(map[string]interface{})

					if int((spinInfo["free_spins"].(float64) + spinInfo["remain_spins"].(float64))) != 0 {
						helper.PrettyLog("success", fmt.Sprintf("%s | Spin Wheel Successfully | Reward: %s %s | Spin Balance: %.0f | Sleep 5s Before Next Spin...", account.Username, reward["reward_value"].(string), reward["reward_type"].(string), (spinInfo["free_spins"].(float64)+spinInfo["remain_spins"].(float64))))
						time.Sleep(5 * time.Second)
					} else {
						helper.PrettyLog("success", fmt.Sprintf("%s | Spin Wheel Successfully | Reward: %s %s | Spin Balance: %.0f", account.Username, reward["reward_value"].(string), reward["reward_type"].(string), (spinInfo["free_spins"].(float64)+spinInfo["remain_spins"].(float64))))

						limit = true
					}
				}
			}
		}
	}

	if isAutoStake {
		stakingInfo := client.getStakingInfo()

		if stake, exits := stakingInfo["stake"].(map[string]interface{}); exits {

			stakingPool := stakingInfo["pool"].(map[string]interface{})

			stakingAmount := fmt.Sprint(stake["available_score"].(float64) * amountStake)

			staking := client.autoStaking(stakingAmount, int(stakingPool["pool_id"].(float64)))

			if staking == nil {
				helper.PrettyLog("error", fmt.Sprintf("%s | Get Staking Info After Auto Stacking Failed...", account.Username))
				return
			}

			helper.PrettyLog("success", fmt.Sprintf("%s | Staking %s Score Successfully  | Current Balance: %.0f | Total Stacked: %s | Realized: %.0f | Pending: %.0f | Total Reward: %.0f | Claimable: %v", account.Username, stakingAmount, staking["available_score"].(float64), staking["stake_score"].(string), staking["realized"].(float64), staking["pending"].(float64), staking["total_reward"].(float64), staking["has_claimables"].(bool)))
		}
	}
}
