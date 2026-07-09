package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type slackPayload struct {
	Text string `json:"text"`
}

func defaultConfigPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "pipe2slack", "webhook_url")
}

func main() {
	configPath := flag.String("config", defaultConfigPath(), "path to a file containing the Slack webhook URL")
	showVersion := flag.Bool("version", false, "print version information and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("pipe2slack %s (commit %s, built %s)\n", version, commit, date)
		return
	}

	if *configPath == "" {
		log.Fatal("No webhook URL file: pass -config, or set $HOME / $XDG_CONFIG_HOME")
	}

	webhookURL, err := loadWebhookURL(*configPath)
	if err != nil {
		log.Fatalf("Error loading webhook URL from %s: %v", *configPath, err)
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Error reading standard input: ", err)
	}

	message := strings.TrimRight(string(input), "\n")
	if message == "" {
		log.Println("No output to send to Slack.")
		return
	}

	if err := sendSlackNotification(webhookURL, message); err != nil {
		log.Fatalf("Error sending Slack notification: %v", err)
	}
	log.Println("Notification sent to Slack successfully.")
}

func loadWebhookURL(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	url := strings.TrimSpace(string(data))
	if url == "" {
		return "", fmt.Errorf("webhook URL is empty")
	}
	return url, nil
}

func sendSlackNotification(webhookURL, message string) error {
	body, err := json.Marshal(slackPayload{Text: message})
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack returned status %s", resp.Status)
	}
	return nil
}
