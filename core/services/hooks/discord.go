package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"response-std/config"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type DiscordHook struct {
	WebhookURL string
	AppName    string
	MinLevel   logrus.Level
}

type DiscordPayload struct {
	Content string  `json:"content,omitempty"`
	Embeds  []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Color       int       `json:"color,omitempty"`
	Fields      []Field   `json:"fields,omitempty"`
	Footer      Footer    `json:"footer,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type Footer struct {
	Text string `json:"text,omitempty"`
}

func NewDiscordHook(webhookURL, appName string, minLevel logrus.Level) *DiscordHook {
	return &DiscordHook{
		WebhookURL: webhookURL,
		AppName:    appName,
		MinLevel:   minLevel,
	}
}

func (hook *DiscordHook) Levels() []logrus.Level {
	var levels []logrus.Level
	for _, level := range logrus.AllLevels {
		if level <= hook.MinLevel {
			levels = append(levels, level)
		}
	}
	return levels
}

func (hook *DiscordHook) Fire(entry *logrus.Entry) error {
	// Skip if webhook URL is empty
	if hook.WebhookURL == "" {
		return nil
	}

	// Create Discord embed
	embed := hook.createEmbed(entry)
	payload := DiscordPayload{
		Embeds: []Embed{embed},
	}

	// Send to Discord
	go hook.sendToDiscord(payload) // Send asynchronously to avoid blocking
	return nil
}

func (hook *DiscordHook) createEmbed(entry *logrus.Entry) Embed {
	// Determine color based on log level
	color := hook.getColorForLevel(entry.Level)

	// Create fields from entry data
	var fields []Field

	// Add basic fields
	if entry.Data != nil {
		for key, value := range entry.Data {
			if key == "error" {
				continue // Handle error separately
			}
			fields = append(fields, Field{
				Name:   strings.Title(strings.ReplaceAll(key, "_", " ")),
				Value:  fmt.Sprintf("%v", value),
				Inline: true,
			})
		}
	}

	// Add error field if exists
	if err, ok := entry.Data["error"]; ok && err != nil {
		fields = append(fields, Field{
			Name:   "Error",
			Value:  fmt.Sprintf("```%v```", err),
			Inline: false,
		})
	}

	embed := Embed{
		Title:       fmt.Sprintf("%s - %s", hook.AppName, strings.ToUpper(entry.Level.String())),
		Description: entry.Message,
		Color:       color,
		Fields:      fields,
		Footer: Footer{
			Text: fmt.Sprintf("Environment: %s", config.ENV.Environment),
		},
		Timestamp: entry.Time,
	}

	return embed
}

func (hook *DiscordHook) getColorForLevel(level logrus.Level) int {
	switch level {
	case logrus.FatalLevel, logrus.PanicLevel:
		return 0xFF0000 // Red
	case logrus.ErrorLevel:
		return 0xFF6B6B // Light Red
	case logrus.WarnLevel:
		return 0xFFB347 // Orange
	case logrus.InfoLevel:
		return 0x4ECDC4 // Teal
	case logrus.DebugLevel:
		return 0x95A5A6 // Gray
	default:
		return 0x3498DB // Blue
	}
}

func (hook *DiscordHook) sendToDiscord(payload DiscordPayload) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(hook.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

// Helper function to parse log level from string
func ParseLogLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal", "critical":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.ErrorLevel
	}
}



// manually send log to discord
func SendDiscordMessage(webhookURL, appName, level, message string, data map[string]interface{}) error {
	color := getColorForStringLevel(level)

	var fields []Field
	for key, value := range data {
		if key == "error" {
			continue
		}
		fields = append(fields, Field{
			Name:   strings.Title(strings.ReplaceAll(key, "_", " ")),
			Value:  fmt.Sprintf("%v", value),
			Inline: true,
		})
	}

	if errVal, ok := data["error"]; ok && errVal != nil {
		fields = append(fields, Field{
			Name:   "Error",
			Value:  fmt.Sprintf("```%v```", errVal),
			Inline: false,
		})
	}

	embed := Embed{
		Title:       fmt.Sprintf("%s - %s", appName, strings.ToUpper(level)),
		Description: message,
		Color:       color,
		Fields:      fields,
		Footer: Footer{
			Text: fmt.Sprintf("Environment: %s", config.ENV.Environment),
		},
		Timestamp: time.Now(),
	}

	payload := DiscordPayload{
		Embeds: []Embed{embed},
	}

	return sendToWebhook(webhookURL, payload)
}

func sendToWebhook(webhookURL string, payload DiscordPayload) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func getColorForStringLevel(level string) int {
	switch strings.ToLower(level) {
	case "fatal", "panic":
		return 0xFF0000
	case "error":
		return 0xFF6B6B
	case "warn", "warning":
		return 0xFFB347
	case "info":
		return 0x4ECDC4
	case "debug":
		return 0x95A5A6
	default:
		return 0x3498DB
	}
}