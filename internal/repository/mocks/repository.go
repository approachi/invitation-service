package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type InvitationRepository struct {
	mock.Mock
}

func (m *InvitationRepository) Use(ctx context.Context, code, email string) error {
	args := m.Called(ctx, code, email)
	return args.Error(0)
}
