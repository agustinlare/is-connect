package main

import (
	"fmt"
	"io"
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

	fmt.Println("Starting is-connect service")

	for {
		if checkInternetConnection() {
			if internetIsLost {
				internetIsLost = false
				internetReturnTime := time.Now()
				fmt.Println("STATUS UP: Internet return at", internetReturnTime.Format("2006-01-02 15:04:05"))
				internetLostDuration := internetReturnTime.Sub(internetLostTime)

				if internetLostDuration.Minutes() >= 1 {
					message := fmt.Sprintf("STATUS UP: Internet was lost by %v minutes", internetLostDuration.Minutes())
					sendDiscordNotification(webhookUrl, message)
				}
			}

		} else {
			if !internetIsLost {
				internetIsLost = true
				internetLostTime = time.Now()
				fmt.Println("STATUS DOWN: Internet was lost at", internetLostTime.Format("2006-01-02 15:04:05"))
			}
		}
		fmt.Println("Internet check", time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(10 * time.Second)
	}
}

func checkInternetConnection() bool {
	website := "https://google.com"
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(website)
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
		log.Println(err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bodyBytes))
}
