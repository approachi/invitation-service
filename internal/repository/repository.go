package repository

import "context"

type InvitationRepository interface {
	Use(ctx context.Context, code, email string) error
}
