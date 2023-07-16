package repository

import "errors"

var (
	ErrFailedSyncStorage   = errors.New("failed to sync emails storage")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrEmailDoesNotExists  = errors.New("email does not exists")
	ErrInvalidEmailAddress = errors.New("invalid email address")
)
