package notifier

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestClientNotifyTo(t *testing.T) {
	customErr := errors.New("custom error")

	tests := []struct {
		name        string
		mockApplier func(m *DialerMock)
		expected    error
	}{
		{
			name: "return error",
			mockApplier: func(m *DialerMock) {
				m.On("DialAndSend", mock.Anything).Return(customErr).Once()
			},
			expected: fmt.Errorf("unexpected error sending mail due to: %w", customErr),
		},
		{
			name: "no error",
			mockApplier: func(m *DialerMock) {
				m.On("DialAndSend", mock.Anything).Return(nil).Once()
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dMock := NewDialerMock(t)

			test.mockApplier(dMock)

			c := client{
				sender: "sender",
				dialer: dMock,
			}

			assert.Equal(t, test.expected, c.NotifyTo(context.TODO(), NotifyToOptions{
				To:      "email",
				Subject: "status",
				Body:    "message",
			}))
		})
	}
}
