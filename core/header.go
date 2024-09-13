package core

import (
	"fmt"
	"net/http"

	"capybara-meme/helper"
)

func setHeader(http *http.Request, referUrl string, authToken string) {

	userAgent, os := generateRandomUserAgent()
	if userAgent == "" || os == "" {
		helper.PrettyLog("error", "Failed Generate Random User Agent")
		return
	}

	header := map[string]string{
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "en-US,en;q=0.9,id;q=0.8",
		"content-type":       "application/json",
		"priority":           "u=1, i",
		"sec-ch-ua":          "\"Chromium\";v=\"128\", \"Not;A=Brand\";v=\"24\", \"Google Chrome\";v=\"128\"",
		"sec-ch-ua-platform": fmt.Sprintf("\"%s\"", os),
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"Referer":            referUrl,
		"Referrer-Policy":    "strict-origin-when-cross-origin",
		"X-Requested-With":   "org.telegram.messenger.web",
		"User-Agent":         userAgent,
	}

	if authToken != "" {
		header["authorization"] = authToken
	}

	for key, value := range header {
		http.Header.Set(key, value)
	}
}
