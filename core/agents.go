package core

import (
	"github.com/mileusna/useragent"

	"capybara-meme/helper"
)

func generateRandomUserAgent() (string, string) {
	userAgents := helper.ReadFileTxt("./core/useragent.txt")
	if userAgents == nil {
		helper.PrettyLog("error", "userAgent data not found")
		return "", ""
	}

	userAgent := userAgents[helper.RandomNumber(0, len(userAgents))]

	os := useragent.Parse(userAgent).OS

	return userAgent, os
}
