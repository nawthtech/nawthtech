package slack

import (
	"context"
	"os"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSlackAPI is a mock implementation of Slack API
type MockSlackAPI struct {
	mock.Mock
}

func (m *MockSlackAPI) PostMessage(channelURL string, options ...slack.MsgOption) (string, string, error) {
	args := m.Called(channelURL, options)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockSlackAPI) UpdateMessage(channelURL, timestamp string, options ...slack.MsgOption) (string, string, string, error) {
	args := m.Called(channelURL, timestamp, options)
	return args.String(0), args.String(1), args.String(2), args.Error(3)
}

func (m *MockSlackAPI) DeleteMessage(channelURL, timestamp string) (string, string, error) {
	args := m.Called(channelURL, timestamp)
	return args.String(0), args.String(1), args.Error(2)
}

func TestNewSlackClient(t *testing.T) {
	tests := []struct {
		name        string
		options     []Option
		wantErr     bool
		errContains string
	}{
		{
			name: "successful creation with token and channel",
			options: []Option{
				WithToken("test-token"),
				WithchannelURL("test-channel"),
				WithAppName("nawthtech"),
				WithEnvironment("test"),
			},
			wantErr: false,
		},
		{
			name: "missing token",
			options: []Option{
				WithchannelURL("test-channel"),
			},
			wantErr:     true,
			errContains: "slack token is required",
		},
		{
			name: "empty token",
			options: []Option{
				WithToken(""),
				WithchannelURL("test-channel"),
			},
			wantErr:     true,
			errContains: "slack token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.options...)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)

				slackClient, ok := client.(*slackClient)
				assert.True(t, ok)
				assert.Equal(t, "nawthtech", slackClient.appName)
				assert.Equal(t, "test", slackClient.env)
			}
		})
	}
}

func TestSlackClient_PushMessage(t *testing.T) {
	mockAPI := new(MockSlackAPI)

	client := &slackClient{
		api:        mockAPI,
		channelURL: "test-channel",
		token:      "test-token",
		appName:    "nawthtech",
		env:        "test",
	}

	t.Run("successful message push", func(t *testing.T) {
		mockAPI.On("PostMessage", "test-channel", mock.Anything).
			Return("test-channel", "1234567890.123456", nil).
			Once()

		channelURL, timestamp, err := client.PushMessage("Hello, world!")

		assert.NoError(t, err)
		assert.Equal(t, "test-channel", channelURL)
		assert.Equal(t, "1234567890.123456", timestamp)
		mockAPI.AssertExpectations(t)
	})

	t.Run("missing channel ID", func(t *testing.T) {
		clientNoChannel := &slackClient{
			api:     mockAPI,
			token:   "test-token",
			appName: "nawthtech",
			env:     "test",
		}

		_, _, err := clientNoChannel.PushMessage("Hello, world!")

		assert.Error(t, err)
		assert.Equal(t, ErrMissingchannelURL, err)
	})

	t.Run("nil client", func(t *testing.T) {
		var nilClient *slackClient
		_, _, err := nilClient.PushMessage("Hello, world!")

		assert.Error(t, err)
		assert.Equal(t, ErrClientNotInitialized, err)
	})

	t.Run("API error", func(t *testing.T) {
		mockAPI.On("PostMessage", "test-channel", mock.Anything).
			Return("", "", assert.AnError).
			Once()

		_, _, err := client.PushMessage("Hello, world!")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to push message")
		mockAPI.AssertExpectations(t)
	})
}

func TestSlackClient_Channel(t *testing.T) {
	baseClient := &slackClient{
		api:        slack.New("test-token"),
		channelURL: "default-channel",
		token:      "test-token",
		appName:    "nawthtech",
		env:        "production",
	}

	t.Run("create channel-specific client", func(t *testing.T) {
		specificClient := baseClient.Channel("specific-channel")

		assert.NotNil(t, specificClient)
		assert.Equal(t, "specific-channel", specificClient.channelURL)
		assert.Equal(t, baseClient.api, specificClient.api)
		assert.Equal(t, baseClient.token, specificClient.token)
		assert.Equal(t, baseClient.appName, specificClient.appName)
		assert.Equal(t, baseClient.env, specificClient.env)
	})

	t.Run("nil base client", func(t *testing.T) {
		var nilClient *slackClient
		specificClient := nilClient.Channel("specific-channel")

		assert.Nil(t, specificClient)
	})
}

func TestSlackClient_PushMessageToChannel(t *testing.T) {
	mockAPI := new(MockSlackAPI)

	client := &slackClient{
		api:        mockAPI,
		channelURL: "default-channel",
		token:      "test-token",
		appName:    "nawthtech",
		env:        "test",
	}

	t.Run("successful message to specific channel", func(t *testing.T) {
		mockAPI.On("PostMessage", "specific-channel", mock.Anything).
			Return("specific-channel", "1234567890.123456", nil).
			Once()

		channelURL, timestamp, err := client.PushMessageToChannel("specific-channel", "Test message")

		assert.NoError(t, err)
		assert.Equal(t, "specific-channel", channelURL)
		assert.Equal(t, "1234567890.123456", timestamp)
		mockAPI.AssertExpectations(t)
	})
}

func TestSlackClient_SendAlert(t *testing.T) {
	mockAPI := new(MockSlackAPI)

	client := &slackClient{
		api:        mockAPI,
		channelURL: "alerts-channel",
		token:      "test-token",
		appName:    "nawthtech",
		env:        "production",
	}

	t.Run("send error alert", func(t *testing.T) {
		mockAPI.On("PostMessage", "alerts-channel", mock.Anything).
			Return("alerts-channel", "1234567890.123456", nil).
			Once()

		channelURL, timestamp, err := client.SendAlert("error", "Database Error", "Failed to connect to database")

		assert.NoError(t, err)
		assert.Equal(t, "alerts-channel", channelURL)
		assert.Equal(t, "1234567890.123456", timestamp)
		mockAPI.AssertExpectations(t)
	})

	t.Run("send success alert", func(t *testing.T) {
		mockAPI.On("PostMessage", "alerts-channel", mock.Anything).
			Return("alerts-channel", "1234567890.123456", nil).
			Once()

		channelURL, timestamp, err := client.SendAlert("success", "Deployment Complete", "Backend deployed successfully")

		assert.NoError(t, err)
		assert.Equal(t, "alerts-channel", channelURL)
		assert.Equal(t, "1234567890.123456", timestamp)
		mockAPI.AssertExpectations(t)
	})
}

func TestSlackClient_SendDeploymentNotification(t *testing.T) {
	mockAPI := new(MockSlackAPI)

	client := &slackClient{
		api:        mockAPI,
		channelURL: "deployments-channel",
		token:      "test-token",
		appName:    "nawthtech",
		env:        "staging",
	}

	t.Run("send successful deployment notification", func(t *testing.T) {
		mockAPI.On("PostMessage", "deployments-channel", mock.Anything).
			Return("deployments-channel", "1234567890.123456", nil).
			Once()

		channelURL, timestamp, err := client.SendDeploymentNotification(
			"backend",
			"v1.2.3",
			"success",
			"abc123def",
			"Fix database connection issue",
		)

		assert.NoError(t, err)
		assert.Equal(t, "deployments-channel", channelURL)
		assert.Equal(t, "1234567890.123456", timestamp)
		mockAPI.AssertExpectations(t)
	})
}

func TestSlackClient_SendErrorNotification(t *testing.T) {
	mockAPI := new(MockSlackAPI)

	client := &slackClient{
		api:        mockAPI,
		channelURL: "errors-channel",
		token:      "test-token",
		appName:    "nawthtech",
		env:        "production",
	}

	t.Run("send error notification with context", func(t *testing.T) {
		mockAPI.On("PostMessage", "errors-channel", mock.Anything).
			Return("errors-channel", "1234567890.123456", nil).
			Once()

		err := assert.AnError
		contextInfo := map[string]string{
			"Service":   "backend",
			"Endpoint":  "/api/users",
			"UserID":    "user-123",
			"RequestID": "req-456",
		}

		channelURL, timestamp, err := client.SendErrorNotification(err, contextInfo)

		assert.NoError(t, err)
		assert.Equal(t, "errors-channel", channelURL)
		assert.Equal(t, "1234567890.123456", timestamp)
		mockAPI.AssertExpectations(t)
	})
}

func TestIntegration(t *testing.T) {
	// Skip integration tests in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// These tests require actual Slack credentials
	// They should only run when SLACK_TOKEN and SLACK_CHANNEL are set
	token := os.Getenv("SLACK_TOKEN")
	channelURL := os.Getenv("SLACK_CHANNEL")

	if token == "" || channelURL == "" {
		t.Skip("Skipping integration test: SLACK_TOKEN or SLACK_CHANNEL not set")
	}

	t.Run("real slack integration", func(t *testing.T) {
		client, err := New(
			WithToken(token),
			WithchannelURL(channelURL),
			WithAppName("nawthtech-test"),
			WithEnvironment("integration-test"),
		)

		assert.NoError(t, err)
		assert.NotNil(t, client)

		// Test sending a simple message
		_, _, err = client.PushMessage("Integration test message from nawthtech")
		assert.NoError(t, err)

		// Test sending an alert
		_, _, err = client.SendAlert("info", "Integration Test", "This is a test alert from integration tests")
		assert.NoError(t, err)
	})
}

func TestClientSingleton(t *testing.T) {
	// Reset default client before test
	defaultClient = nil

	t.Run("client not initialized", func(t *testing.T) {
		client := Client()
		assert.Nil(t, client)
	})

	t.Run("client initialized", func(t *testing.T) {
		mockAPI := new(MockSlackAPI)
		defaultClient = &slackClient{
			api:        mockAPI,
			channelURL: "test-channel",
			token:      "test-token",
		}

		client := Client()
		assert.NotNil(t, client)
	})
}
