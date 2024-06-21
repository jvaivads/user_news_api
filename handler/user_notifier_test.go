package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"user_news_api/handler/mocks"
	"user_news_api/services"
)

func TestHandleNotifyUser(t *testing.T) {
	tests := []struct {
		name           string
		payload        NotifyUserRequestPayload
		setupMocks     func(service *mocks.UserNotifier)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid request",
			payload: NotifyUserRequestPayload{
				UserEmail:   "test@example.com",
				MessageType: "welcome",
			},
			setupMocks: func(service *mocks.UserNotifier) {
				service.On("Notify", mock.Anything, "test@example.com", "welcome").Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name: "Invalid email format",
			payload: NotifyUserRequestPayload{
				UserEmail:   "invalid-email",
				MessageType: "welcome",
			},
			setupMocks:     func(service *mocks.UserNotifier) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "request validation fails due to: Key: 'NotifyUserRequestPayload.UserEmail' Error:Field validation for 'UserEmail' failed on the 'email' tag",
		},
		{
			name: "Empty message type",
			payload: NotifyUserRequestPayload{
				UserEmail:   "test@example.com",
				MessageType: "",
			},
			setupMocks:     func(service *mocks.UserNotifier) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "request validation fails due to: Key: 'NotifyUserRequestPayload.MessageType' Error:Field validation for 'MessageType' failed on the 'required' tag",
		},
		{
			name: "Service limit exceeded error",
			payload: NotifyUserRequestPayload{
				UserEmail:   "test@example.com",
				MessageType: "welcome",
			},
			setupMocks: func(service *mocks.UserNotifier) {
				service.On("Notify", mock.Anything, "test@example.com", "welcome").
					Return(fmt.Errorf(
						"%w: rate limit reached for user", services.ErrLimitExceeded)).Once()
			},
			expectedStatus: http.StatusTooManyRequests,
			expectedBody:   "too many requests",
		},
		{
			name: "Service internal error",
			payload: NotifyUserRequestPayload{
				UserEmail:   "test@example.com",
				MessageType: "welcome",
			},
			setupMocks: func(service *mocks.UserNotifier) {
				service.On("Notify", mock.Anything, "test@example.com", "welcome").
					Return(errors.New("internal error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewUserNotifier(t)

			tt.setupMocks(mockService)

			controller := &UserController{service: mockService}

			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			controller.handleNotifyUser(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedBody != "" {
				body, err = io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Contains(t, string(body), tt.expectedBody)
			}
		})
	}
}
