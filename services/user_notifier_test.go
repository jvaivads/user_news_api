package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"user_news_api/notifier"
	"user_news_api/ratelimiter"
	"user_news_api/services/mocks"

	"github.com/stretchr/testify/assert"
)

func TestUserNotifier_Notify(t *testing.T) {
	ctx := context.Background()
	userMail := "user@example.com"
	messageType := ratelimiter.NewsType

	tests := []struct {
		name          string
		applyMocks    func(*mocks.Limiter, *mocks.Notifier)
		expectedError error
	}{
		{
			name: "Success",
			applyMocks: func(ml *mocks.Limiter, mn *mocks.Notifier) {
				ml.On("Reached", ctx, userMail, messageType).Return(false, nil).Once()
				mn.On("NotifyTo", ctx, notifier.NotifyToOptions{
					To:      userMail,
					Subject: "Notification",
					Body:    toHTML(messageType),
				}).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Limiter Error",
			applyMocks: func(ml *mocks.Limiter, mn *mocks.Notifier) {
				ml.On("Reached", ctx, userMail, messageType).Return(false, errors.New("limiter error")).Once()
			},
			expectedError: fmt.Errorf("limiter error for user %s: %w", userMail, errors.New("limiter error")),
		},
		{
			name: "Rate Limit Exceeded",
			applyMocks: func(ml *mocks.Limiter, mn *mocks.Notifier) {
				ml.On("Reached", ctx, userMail, messageType).Return(true, nil).Once()
			},
			expectedError: fmt.Errorf("%w: rate limit reached for user %s and message type %s", ErrLimitExceeded, userMail, messageType),
		},
		{
			name: "Notifier Error",
			applyMocks: func(ml *mocks.Limiter, mn *mocks.Notifier) {
				ml.On("Reached", ctx, userMail, messageType).Return(false, nil).Once()
				mn.On("NotifyTo", ctx, notifier.NotifyToOptions{
					To:      userMail,
					Subject: "Notification",
					Body:    toHTML(messageType),
				}).Return(errors.New("notifier error")).Once()
			},
			expectedError: fmt.Errorf("notifier error for user %s: %w", userMail, errors.New("notifier error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLimiter := mocks.NewLimiter(t)
			mockNotifier := mocks.NewNotifier(t)

			tt.applyMocks(mockLimiter, mockNotifier)

			serv := UserNotifier{
				limiter:  mockLimiter,
				notifier: mockNotifier,
			}

			err := serv.Notify(ctx, userMail, messageType)

			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestToHTML(t *testing.T) {
	tests := []struct {
		messageType string
		expected    string
	}{
		{
			messageType: ratelimiter.NewsType,
			expected: `
<!DOCTYPE html>
<html>
<head>
    <style>
        .centered {
            text-align: center;
            color: green;
            font-size: 48px;
        }
    </style>
</head>
<body>
    <div class="centered">News</div>
</body>
</html>
`,
		},
		{
			messageType: ratelimiter.StatusType,
			expected: `
<!DOCTYPE html>
<html>
<head>
    <style>
        .centered {
            text-align: center;
            color: red;
            font-size: 48px;
        }
    </style>
</head>
<body>
    <div class="centered">Status</div>
</body>
</html>
`,
		},
		{
			messageType: ratelimiter.MarketingType,
			expected: `
<!DOCTYPE html>
<html>
<head>
    <style>
        .centered {
            text-align: center;
            color: yellow;
            font-size: 48px;
        }
    </style>
</head>
<body>
    <div class="centered">Marketing</div>
</body>
</html>
`,
		},
		{
			messageType: "other",
			expected: `
<!DOCTYPE html>
<html>
<head>
    <style>
        .centered {
            text-align: center;
            color: black;
            font-size: 48px;
        }
    </style>
</head>
<body>
    <div class="centered">other</div>
</body>
</html>
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.messageType, func(t *testing.T) {
			result := toHTML(tt.messageType)

			assert.Equal(t, tt.expected, result)
		})
	}
}
