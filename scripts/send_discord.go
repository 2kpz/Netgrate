package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func DiscordRun() string {
	webhookURL := ""

	payload := map[string]string{
		"content": "Launching Chat Application!",
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Sprintf("[Error preparing payload: %v]", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Sprintf("[Error sending request to Discord: %v]", err)
	}
	defer resp.Body.Close()

	return "Chat Application Launched!"
}

func main() {
	fmt.Print(DiscordRun())
}
