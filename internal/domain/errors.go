package domain

import "errors"

var (
	ErrInvitationNotFound  = errors.New("invitation not found")
	ErrInvitationExhausted = errors.New("invitation exhausted")
	ErrAlreadyRegistered   = errors.New("already registered with this invitation")
	ErrInvalidEmail        = errors.New("invalid email")
)

func IsClientError(err error) bool {
	return errors.Is(err, ErrInvalidEmail) ||
		errors.Is(err, ErrInvitationNotFound) ||
		errors.Is(err, ErrInvitationExhausted) ||
		errors.Is(err, ErrAlreadyRegistered)
}
