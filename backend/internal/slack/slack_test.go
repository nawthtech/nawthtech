package slack

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestNewSlackClient_Validation(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
		wantErr bool
	}{
		{
			name:    "missing token",
			options: []Option{},
			wantErr: true,
		},
		{
			name: "valid config",
			options: []Option{
				WithToken("test-token"),
				WithChannelURL("test-channel"),
			},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.options...)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

// اختبارات بسيطة بدون mock معقد
func TestSlackClient_Basic(t *testing.T) {
	client := &SlackClient{
		channelID: "test-channel",
		token:     "test-token",
	}
	
	assert.Equal(t, "test-channel", client.channelID)
	assert.Equal(t, "test-token", client.token)
}