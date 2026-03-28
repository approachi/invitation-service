package domain

import "errors"

var (
	ErrInvitationNotFound  = errors.New("invitation not found")
	ErrInvitationExhausted = errors.New("invitation exhausted")
	ErrAlreadyRegistered   = errors.New("already registered with this invitation")
	ErrInvalidEmail        = errors.New("invalid email")
)
