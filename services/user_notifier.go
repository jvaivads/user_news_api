package services

import (
	"context"
	"errors"
	"fmt"
	"user_news_api/notifier"
	"user_news_api/ratelimiter"
)

var (
	ErrLimitExceeded = errors.New("limit exceeded")
)

// Limiter is an abstraction for ratelimiter.LimiterPool making it mockeable
type Limiter interface {
	Reached(context.Context, string, string) (bool, error)
}

// Notifier is an abstraction for notifier.Client making it mockeable
type Notifier interface {
	NotifyTo(context.Context, notifier.NotifyToOptions) error
}

func NewUserNotifier(limiter Limiter, notifier Notifier) UserNotifierService {
	return UserNotifierService{
		limiter:  limiter,
		notifier: notifier,
	}
}

type UserNotifierService struct {
	limiter  Limiter
	notifier Notifier
}

func (serv UserNotifierService) Notify(ctx context.Context, userMail string, messageType string) error {
	reached, err := serv.limiter.Reached(ctx, userMail, messageType)
	if err != nil {
		return fmt.Errorf("limiter error for user %s: %w", userMail, err)
	}

	if reached {
		return fmt.Errorf(
			"%w: rate limit reached for user %s and message type %s", ErrLimitExceeded, userMail, messageType)
	}

	err = serv.notifier.NotifyTo(ctx, notifier.NotifyToOptions{
		To:      userMail,
		Subject: "Notification",
		Body:    toHTML(messageType),
	})
	if err != nil {
		return fmt.Errorf("notifier error for user %s: %w", userMail, err)
	}

	return nil
}

// toHTML is only string formatter, and it is used by the service like a decorator.
func toHTML(messageType string) string {
	color := "black"
	switch messageType {
	case ratelimiter.NewsType:
		color = "green"
	case ratelimiter.StatusType:
		color = "red"
	case ratelimiter.MarketingType:
		color = "yellow"
	}

	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <style>
        .centered {
            text-align: center;
            color: %s;
            font-size: 48px;
        }
    </style>
</head>
<body>
    <div class="centered">%s</div>
</body>
</html>
`
	return fmt.Sprintf(htmlTemplate, color, messageType)
}
