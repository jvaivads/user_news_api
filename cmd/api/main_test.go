package main

import (
	"os"
	"testing"
	"user_news_api/notifier"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNotifierOptions(t *testing.T) {
	tests := []struct {
		name         string
		envVars      map[string]string
		expectedOpts notifier.Options
		expectPanic  bool
		panicMessage string
	}{
		{
			name: "All environment variables set correctly",
			envVars: map[string]string{
				"NOTIFIER_HOST":     "smtp.example.com",
				"NOTIFIER_PORT":     "587",
				"NOTIFIER_SENDER":   "user@example.com",
				"NOTIFIER_PASSWORD": "password",
			},
			expectedOpts: notifier.Options{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "user@example.com",
				Password: "password",
			},
			expectPanic: false,
		},
		{
			name: "Notifier host is empty",
			envVars: map[string]string{
				"NOTIFIER_PORT":     "587",
				"NOTIFIER_SENDER":   "user@example.com",
				"NOTIFIER_PASSWORD": "password",
			},
			expectPanic:  true,
			panicMessage: "notifier host is empty",
		},
		{
			name: "Notifier port is empty",
			envVars: map[string]string{
				"NOTIFIER_HOST":     "smtp.example.com",
				"NOTIFIER_SENDER":   "user@example.com",
				"NOTIFIER_PASSWORD": "password",
			},
			expectPanic:  true,
			panicMessage: "notifier port is empty",
		},
		{
			name: "Notifier port is not a number",
			envVars: map[string]string{
				"NOTIFIER_HOST":     "smtp.example.com",
				"NOTIFIER_PORT":     "not a number",
				"NOTIFIER_SENDER":   "user@example.com",
				"NOTIFIER_PASSWORD": "password",
			},
			expectPanic:  true,
			panicMessage: "notifier port is not a number",
		},
		{
			name: "Notifier sender is empty",
			envVars: map[string]string{
				"NOTIFIER_HOST":     "smtp.example.com",
				"NOTIFIER_PORT":     "587",
				"NOTIFIER_PASSWORD": "password",
			},
			expectPanic:  true,
			panicMessage: "notifier sender is empty",
		},
		{
			name: "Notifier password is empty",
			envVars: map[string]string{
				"NOTIFIER_HOST":   "smtp.example.com",
				"NOTIFIER_PORT":   "587",
				"NOTIFIER_SENDER": "user@example.com",
			},
			expectPanic:  true,
			panicMessage: "notifier password is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				require.NoError(t, os.Setenv(key, value))
			}

			defer func() {
				for key := range tt.envVars {
					require.NoError(t, os.Unsetenv(key))
				}
			}()

			if tt.expectPanic {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, tt.panicMessage, r)
					} else {
						t.Errorf("Expected panic with message: %s", tt.panicMessage)
					}
				}()
			}

			assert.Equal(t, tt.expectedOpts, getNotifierOptions())
		})
	}
}

func TestGetRedisOptions(t *testing.T) {
	tests := []struct {
		name         string
		envVars      map[string]string
		expectedOpts *redis.Options
		expectPanic  bool
		panicMessage string
	}{
		{
			name: "Redis address is empty",
			envVars: map[string]string{
				"REDIS_ADDRESS": "",
			},
			expectPanic:  true,
			panicMessage: "redis address is empty",
		},
		{
			name: "OK",
			envVars: map[string]string{
				"REDIS_ADDRESS":  "address",
				"REDIS_PASSWORD": "pass",
			},
			expectPanic:  false,
			panicMessage: "",
			expectedOpts: &redis.Options{
				Addr:     "address",
				Password: "pass",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				require.NoError(t, os.Setenv(key, value))
			}

			defer func() {
				for key := range tt.envVars {
					require.NoError(t, os.Unsetenv(key))
				}
			}()

			if tt.expectPanic {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, tt.panicMessage, r)
					} else {
						t.Errorf("Expected panic with message: %s", tt.panicMessage)
					}
				}()
			}

			assert.Equal(t, tt.expectedOpts, getRedisOptions())
		})
	}
}
