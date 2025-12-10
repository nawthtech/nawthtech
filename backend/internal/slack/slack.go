package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/slack-go/slack"
)

// SlackClient defines the interface for Slack operations
type SlackClient interface {
	PushMessage(text string) (string, string, error)
	PushMessageWithAttachments(text string, attachments []slack.Attachment) (string, string, error)
	PushMessageToChannel(channel, text string) (string, string, error)
	UpdateMessage(channelURL, timestamp, text string) (string, string, string, error)
	DeleteMessage(channelURL, timestamp string) (string, error)
 SendAlert(alertType, title, message string) (string, string, error) 
 SendDeploymentNotification(service, version, status, commitHash, commitMessage string) (string, string, error)
 SendErrorNotification (err error, contextInfo map[stringlstring) (string, string, error)
}

type slackClient struct {
	api        *slack.Client
	channelURL string
	token      string
	appName    string
	env        string // production, staging, development
}

type Option func(*slackClient) error

var (
	// Default client instance
	defaultClient *slackClient

	// Error messages
	ErrClientNotInitialized = fmt.Errorf("slack client not initialized")
	ErrMissingToken         = fmt.Errorf("slack token is required")
	ErrMissingchannelURL    = fmt.Errorf("slack channel ID is required")
)

// New creates a new Slack client with the provided options
func New(options ...Option) (SlackClient, error) {
	client := &slackClient{}

	for _, opt := range options {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	// Validate required fields
	if client.token == "" {
		return nil, ErrMissingToken
	}

	// Initialize Slack API
	client.api = slack.New(client.token)

	// Set default values if not provided
	if client.appName == "" {
		client.appName = "nawthtech"
	}
	if client.env == "" {
		client.env = "development"
	}

	// Store as default client
	defaultClient = client

	return client, nil
}

// Init initializes the default client with options
func Init(options ...Option) error {
	_, err := New(options...)
	return err
}

// Client returns the default client instance
func Client() SlackClient { // غيرت نوع الإرجاع لـ SlackClient
	if defaultClient == nil {
		log.Println("[WARN] Slack client not initialized, returning nil")
		return nil
	}
	return defaultClient
}

// WithToken sets the Slack bot token
func WithToken(token string) Option {
	return func(c *slackClient) error {
		c.token = token
		return nil
	}
}

// WithchannelURL sets the default channel ID
func WithchannelURL(channelURL string) Option {
	return func(c *slackClient) error {
		c.channelURL = channelURL
		return nil
	}
}

// WithAppName sets the application name for logging
func WithAppName(appName string) Option {
	return func(c *slackClient) error {
		c.appName = appName
		return nil
	}
}

// WithEnvironment sets the environment (production, staging, development)
func WithEnvironment(env string) Option {
	return func(c *slackClient) error {
		c.env = env
		return nil
	}
}

// Channel returns a new client instance with the specified channel
func (c *slackClient) Channel(channelURL string) SlackClient { // غيرت نوع الإرجاع
	if c == nil {
		return nil
	}

	return &slackClient{
		api:        c.api,
		token:      c.token,
		channelURL: channelURL,
		appName:    c.appName,
		env:        c.env,
	}
}

// PushMessage sends a message to the default channel
func (c *slackClient) PushMessage(text string) (string, string, error) {
	if c == nil || c.api == nil {
		return "", "", ErrClientNotInitialized
	}

	if c.channelURL == "" {
		return "", "", ErrMissingchannelURL
	}

	// Add environment and app name prefix
	formattedText := fmt.Sprintf("[%s:%s] %s", c.env, c.appName, text)

	channelURL, timestamp, err := c.api.PostMessage(
		c.channelURL,
		slack.MsgOptionText(formattedText, false),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		return "", "", fmt.Errorf("failed to push message: %w", err)
	}

	log.Printf("Message sent to channel %s at %s", channelURL, timestamp)
	return channelURL, timestamp, nil
}

// PushMessageWithAttachments sends a message with attachments
func (c *slackClient) PushMessageWithAttachments(text string, attachments []slack.Attachment) (string, string, error) {
	if c == nil || c.api == nil {
		return "", "", ErrClientNotInitialized
	}

	if c.channelURL == "" {
		return "", "", ErrMissingchannelURL
	}

	formattedText := fmt.Sprintf("[%s:%s] %s", c.env, c.appName, text)

	channelURL, timestamp, err := c.api.PostMessage(
		c.channelURL,
		slack.MsgOptionText(formattedText, false),
		slack.MsgOptionAttachments(attachments...),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		return "", "", fmt.Errorf("failed to push message with attachments: %w", err)
	}

	return channelURL, timestamp, nil
}

// PushMessageToChannel sends a message to a specific channel
func (c *slackClient) PushMessageToChannel(channel, text string) (string, string, error) {
	if c == nil || c.api == nil {
		return "", "", ErrClientNotInitialized
	}

	formattedText := fmt.Sprintf("[%s:%s] %s", c.env, c.appName, text)

	channelURL, timestamp, err := c.api.PostMessage(
		channel,
		slack.MsgOptionText(formattedText, false),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		return "", "", fmt.Errorf("failed to push message to channel %s: %w", channel, err)
	}

	return channelURL, timestamp, nil
}

// UpdateMessage updates an existing message
func (c *slackClient) UpdateMessage(channelURL, timestamp, text string) (string, string, string, error) {
	if c == nil || c.api == nil {
		return "", "", "", ErrClientNotInitialized
	}

	formattedText := fmt.Sprintf("[%s:%s] %s", c.env, c.appName, text)

	newTimestamp, _, _, err := c.api.UpdateMessage(
		channelURL,
		timestamp,
		slack.MsgOptionText(formattedText, false),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		return "", "", "", fmt.Errorf("failed to update message: %w", err)
	}

	return channelURL, timestamp, newTimestamp, nil
}

// DeleteMessage deletes a message - التصحيح الرئيسي هنا
func (c *slackClient) DeleteMessage(channelURL, timestamp string) (string, error) {
	if c == nil || c.api == nil {
		return "", ErrClientNotInitialized
	}

	// وفقاً لمكتبة slack-go/slack، DeleteMessage ترجع error فقط
	// إذا كانت ترجع timestamp، استخدم:
	// _, _, err := c.api.DeleteMessage(channelURL, timestamp)
	// return "deleted", err
	
	// أو إذا كانت ترجع string:
	// response, err := c.api.DeleteMessage(channelURL, timestamp)
	// return response, err

	// بناءً على الخطأ: "c.api.DeleteMessage returns 3 values"
	// إذن المكتبة ترجع 3 قيم
	_, _, err := c.api.DeleteMessage(channelURL, timestamp)
	if err != nil {
		return "", fmt.Errorf("failed to delete message: %w", err)
	}

	return "deleted", nil
}

// SendAlert sends an alert message with alert formatting
func (c *slackClient) SendAlert(alertType, title, message string) (string, string, error) {
	var emoji string
	var color string

	switch alertType {
	case "error":
		emoji = "🚨"
		color = "#FF0000"
	case "warning":
		emoji = "⚠️"
		color = "#FFA500"
	case "info":
		emoji = "ℹ️"
		color = "#3498DB"
	case "success":
		emoji = "✅"
		color = "#2ECC71"
	default:
		emoji = "📢"
		color = "#9B59B6"
	}

	formattedTitle := fmt.Sprintf("%s [%s:%s] %s", emoji, c.env, c.appName, title)

	attachment := slack.Attachment{
		Color:      color,
		Title:      formattedTitle,
		Text:       message,
		MarkdownIn: []string{"text", "fields"},
		Footer:     "nawthtech",
		Ts:         json.Number(fmt.Sprintf("%d", time.Now().Unix())),
	}

	return c.PushMessageWithAttachments("", []slack.Attachment{attachment})
}

// SendDeploymentNotification sends a deployment notification
func (c *slackClient) SendDeploymentNotification(service, version, status, commitHash, commitMessage string) (string, string, error) {
	title := fmt.Sprintf("Deployment %s - %s", status, service)

	fields := []slack.AttachmentField{
		{
			Title: "Service",
			Value: service,
			Short: true,
		},
		{
			Title: "Version",
			Value: version,
			Short: true,
		},
		{
			Title: "Environment",
			Value: c.env,
			Short: true,
		},
		{
			Title: "Commit",
			Value: commitHash,
			Short: true,
		},
	}

	if commitMessage != "" {
		fields = append(fields, slack.AttachmentField{
			Title: "Commit Message",
			Value: commitMessage,
			Short: false,
		})
	}

	var color string
	switch status {
	case "success":
		color = "#2ECC71"
	case "failed":
		color = "#E74C3C"
	case "started":
		color = "#3498DB"
	default:
		color = "#9B59B6"
	}

	attachment := slack.Attachment{
		Color:  color,
		Title:  title,
		Fields: fields,
		Footer: "Railway Deployment",
		Ts:     json.Number(fmt.Sprintf("%d", time.Now().Unix())),
	}

	return c.PushMessageWithAttachments("", []slack.Attachment{attachment})
}

// SendErrorNotification sends an error notification
func (c *slackClient) SendErrorNotification(err error, contextInfo map[string]string) (string, string, error) {
	title := fmt.Sprintf("Error in %s", c.appName)

	fields := []slack.AttachmentField{
		{
			Title: "Error Message",
			Value: err.Error(),
			Short: false,
		},
		{
			Title: "Environment",
			Value: c.env,
			Short: true,
		},
		{
			Title: "Timestamp",
			Value: time.Now().Format(time.RFC3339),
			Short: true,
		},
	}

	// Add context information
	for key, value := range contextInfo {
		fields = append(fields, slack.AttachmentField{
			Title: key,
			Value: value,
			Short: true,
		})
	}

	attachment := slack.Attachment{
		Color:  "#E74C3C",
		Title:  title,
		Fields: fields,
		Footer: "Error Notification",
		Ts:     json.Number(fmt.Sprintf("%d", time.Now().Unix())),
	}

	return c.PushMessageWithAttachments("", []slack.Attachment{attachment})
}

// Helper Functions

// DefaultPushMessage uses the default client to push a message
func DefaultPushMessage(text string) (string, string, error) {
	client := Client()
	if client == nil {
		return "", "", ErrClientNotInitialized
	}
	return client.PushMessage(text)
}

// DefaultSendAlert uses the default client to send an alert
func DefaultSendAlert(alertType, title, message string) (string, string, error) {
	client := Client()
	if client == nil {
		return "", "", ErrClientNotInitialized
	}
	return client.SendAlert(alertType, title, message)
}