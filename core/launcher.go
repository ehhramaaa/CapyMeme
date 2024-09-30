package core

import (
	"CapybaraMeme/tools"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gookit/config/v2"
)

func (account *Account) worker(wg *sync.WaitGroup, semaphore *chan struct{}, totalPointsChan *chan int, index int, query string, proxyList []string, walletList []string) {
	defer wg.Done()
	*semaphore <- struct{}{}

	var points int
	var proxy string

	if len(proxyList) > 0 {
		proxy = proxyList[index%len(proxyList)]
	}

	tools.Logger("info", fmt.Sprintf("| %s | Starting Bot...", account.username))

	setDns(&net.Dialer{})

	client := Client{
		account: *account,
		proxy:   proxy,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}

	if len(client.proxy) > 0 {
		err := client.setProxy()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to set proxy: %v", account.username, err))
		} else {
			tools.Logger("success", fmt.Sprintf("| %s | Proxy Successfully Set...", account.username))
		}
	}

	infoIp, err := client.checkIp()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed to check ip: %v", err))
	}

	if infoIp != nil {
		tools.Logger("success", fmt.Sprintf("| %s | Ip: %s | City: %s | Country: %s | Provider: %s", account.username, infoIp["ip"].(string), infoIp["city"].(string), infoIp["country"].(string), infoIp["org"].(string)))
	}

	points = client.autoCompleteTask()

	*totalPointsChan <- points

	<-*semaphore
}

func LaunchBot() {
	queryPath := "configs/query.txt"
	proxyPath := "configs/proxy.txt"
	maxThread := config.Int("MAX_THREAD")
	isUseProxy := config.Bool("USE_PROXY")

	if config.Int("STAKING_PERCENTAGE") > 100 {
		tools.Logger("error", "Staking Amount Percentage Must Be Less Than 100%")
	}

	queryData, err := tools.ReadFileTxt(queryPath)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Query Data Not Found: %s", err))
		return
	}

	tools.Logger("info", fmt.Sprintf("%v Query Data Detected", len(queryData)))

	var wg sync.WaitGroup
	var semaphore chan struct{}
	var proxyList, walletList []string

	if isUseProxy {
		proxyList, err = tools.ReadFileTxt(proxyPath)
		if err != nil {
			tools.Logger("error", fmt.Sprintf("Proxy Data Not Found: %s", err))
		}

		tools.Logger("info", fmt.Sprintf("%v Proxy Detected", len(proxyList)))
	}

	totalPointsChan := make(chan int, len(queryData))

	if maxThread > len(queryData) {
		semaphore = make(chan struct{}, len(queryData))
	} else {
		semaphore = make(chan struct{}, maxThread)
	}

	for {
		for index, query := range queryData {
			wg.Add(1)
			account := &Account{
				queryData: query,
			}

			account.parsingQueryData()

			go account.worker(&wg, &semaphore, &totalPointsChan, index, query, proxyList, walletList)
		}
		wg.Wait()
		close(totalPointsChan)

		var totalPoints int

		for points := range totalPointsChan {
			totalPoints += points
		}

		tools.Logger("success", fmt.Sprintf("Total Points All Account: %v", totalPoints))

		randomSleep := tools.RandomNumber(config.Int("RANDOM_SLEEP.MIN"), config.Int("RANDOM_SLEEP.MAX"))

		tools.Logger("info", fmt.Sprintf("Launch Bot Finished | Sleep %vs Before Next Lap...", randomSleep))

		time.Sleep(time.Duration(randomSleep) * time.Second)
	}
}
