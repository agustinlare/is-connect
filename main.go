package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	internetLostTime := time.Now()
	internetIsLost := false
	webhookUrl := os.Getenv("WEBHOOK_URL")

	if webhookUrl == "" {
		panic("WEBHOOK_URL environment variable is not set or is empty")
	}

	for {
		if checkInternetConnection() {
			if internetIsLost {
				internetIsLost = false
				internetReturnTime := time.Now()
				internetLostDuration := internetReturnTime.Sub(internetLostTime)
				if internetLostDuration.Minutes() >= 1 {
					message := fmt.Sprintf("Internet connection is back after %v minutes", internetLostDuration.Minutes())
					sendDiscordNotification(webhookUrl, message)
				}
			}
		} else {
			if !internetIsLost {
				internetIsLost = true
				internetLostTime = time.Now()
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func checkInternetConnection() bool {
	resp, err := http.Get("https://www.google.com")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

func sendDiscordNotification(webhookUrl, message string) {
	formData := url.Values{
		"content": {message},
	}
	resp, err := http.PostForm(webhookUrl, formData)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
