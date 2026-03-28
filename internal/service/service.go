package service

import "context"

type InvitationService interface {
	UseInvitation(ctx context.Context, code, email string) error
}
