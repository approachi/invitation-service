package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type InvitationService struct {
	mock.Mock
}

func (m *InvitationService) UseInvitation(ctx context.Context, code, email string) error {
	args := m.Called(ctx, code, email)
	return args.Error(0)
}
