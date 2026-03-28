package service_test

import (
	"context"
	"testing"

	"invitation-service/internal/domain"
	"invitation-service/internal/repository/mocks"
	"invitation-service/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseInvitation_ShouldSucceed_WhenEmailAndCodeAreValid(t *testing.T) {
	repo := new(mocks.InvitationRepository)
	repo.On("Use", mock.Anything, "twitter-reg1", "user@example.com").Return(nil)

	svc := service.NewInvitationService(repo)
	err := svc.UseInvitation(context.Background(), "twitter-reg1", "user@example.com")

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUseInvitation_ShouldReturnInvalidEmail_WhenEmailIsInvalid(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"empty email", ""},
		{"missing domain", "notanemail"},
		{"name format", "John <john@example.com>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.InvitationRepository)

			svc := service.NewInvitationService(repo)
			err := svc.UseInvitation(context.Background(), "twitter-reg1", tt.email)

			assert.ErrorIs(t, err, domain.ErrInvalidEmail)
			repo.AssertNotCalled(t, "Use")
		})
	}
}

func TestUseInvitation_ShouldReturnNotFound_WhenCodeIsUnknown(t *testing.T) {
	repo := new(mocks.InvitationRepository)
	repo.On("Use", mock.Anything, "unknown-code", "user@example.com").Return(domain.ErrInvitationNotFound)

	svc := service.NewInvitationService(repo)
	err := svc.UseInvitation(context.Background(), "unknown-code", "user@example.com")

	assert.ErrorIs(t, err, domain.ErrInvitationNotFound)
	repo.AssertExpectations(t)
}

func TestUseInvitation_ShouldReturnExhausted_WhenMaxUsesReached(t *testing.T) {
	repo := new(mocks.InvitationRepository)
	repo.On("Use", mock.Anything, "twitter-reg1", "user@example.com").Return(domain.ErrInvitationExhausted)

	svc := service.NewInvitationService(repo)
	err := svc.UseInvitation(context.Background(), "twitter-reg1", "user@example.com")

	assert.ErrorIs(t, err, domain.ErrInvitationExhausted)
	repo.AssertExpectations(t)
}

func TestUseInvitation_ShouldReturnAlreadyRegistered_WhenEmailAlreadyUsed(t *testing.T) {
	repo := new(mocks.InvitationRepository)
	repo.On("Use", mock.Anything, "twitter-reg1", "user@example.com").Return(domain.ErrAlreadyRegistered)

	svc := service.NewInvitationService(repo)
	err := svc.UseInvitation(context.Background(), "twitter-reg1", "user@example.com")

	assert.ErrorIs(t, err, domain.ErrAlreadyRegistered)
	repo.AssertExpectations(t)
}
