package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"invitation-service/internal/domain"
	"invitation-service/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseInvitation_ShouldReturn200_WhenRequestIsValid(t *testing.T) {
	svc := new(mocks.InvitationService)
	svc.On("UseInvitation", mock.Anything, "twitter-reg1", "user@example.com").Return(nil)

	body := `{"Email":"user@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/invitations/twitter-reg1", strings.NewReader(body))
	rec := httptest.NewRecorder()

	newTestRouter(svc).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	svc.AssertExpectations(t)
}

func TestUseInvitation_ShouldReturn400_WhenBodyIsInvalidJSON(t *testing.T) {
	svc := new(mocks.InvitationService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/invitations/twitter-reg1", strings.NewReader("not-json"))
	rec := httptest.NewRecorder()

	newTestRouter(svc).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	svc.AssertNotCalled(t, "UseInvitation")
}

func TestUseInvitation_ShouldReturn400_WhenServiceReturnsClientError(t *testing.T) {
	tests := []struct {
		name   string
		svcErr error
	}{
		{"invalid email", domain.ErrInvalidEmail},
		{"not found", domain.ErrInvitationNotFound},
		{"exhausted", domain.ErrInvitationExhausted},
		{"already registered", domain.ErrAlreadyRegistered},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.InvitationService)
			svc.On("UseInvitation", mock.Anything, "twitter-reg1", "user@example.com").Return(tt.svcErr)

			body := `{"Email":"user@example.com"}`
			req := httptest.NewRequest(http.MethodPost, "/api/v1/invitations/twitter-reg1", strings.NewReader(body))
			rec := httptest.NewRecorder()

			newTestRouter(svc).ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			svc.AssertExpectations(t)
		})
	}
}

func TestUseInvitation_ShouldReturn500_WhenServiceReturnsUnexpectedError(t *testing.T) {
	svc := new(mocks.InvitationService)
	svc.On("UseInvitation", mock.Anything, "twitter-reg1", "user@example.com").Return(assert.AnError)

	body := `{"Email":"user@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/invitations/twitter-reg1", strings.NewReader(body))
	rec := httptest.NewRecorder()

	newTestRouter(svc).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	svc.AssertExpectations(t)
}
