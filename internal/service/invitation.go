package service

import (
	"context"
	"fmt"
	"net/mail"

	"invitation-service/internal/domain"
	"invitation-service/internal/repository"
)

type invitationService struct {
	repo repository.InvitationRepository
}

func NewInvitationService(repo repository.InvitationRepository) InvitationService {
	return &invitationService{repo: repo}
}

func (s *invitationService) UseInvitation(ctx context.Context, code, email string) error {
	if err := validateEmail(email); err != nil {
		return err
	}
	return s.repo.Use(ctx, code, email)
}

func validateEmail(email string) error {
	addr, err := mail.ParseAddress(email)
	if err == nil && addr.Address == email {
		return nil
	}
	return fmt.Errorf("%w: %s", domain.ErrInvalidEmail, email)
}
