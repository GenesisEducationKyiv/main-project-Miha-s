package emailsStorage

import "errors"

var (
	ErrFailedSyncStorage   = errors.New("failed to sync emails storage")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidEmailAddress = errors.New("invalid email address")
)
