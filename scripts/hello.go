package main

import (
	"fmt"
	"io"
	"net/http"
)

func GetIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve IP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(ip), nil
}

func main() {
	ip, err := GetIP()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Print(ip) // Print the IP address
}
