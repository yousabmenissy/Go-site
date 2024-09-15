package internal

import "errors"

var (
	ErrDuplicateEmail     = errors.New("account with this email already exists")
	ErrNoRecord           = errors.New("the requested record was not found")
	ErrEditConflict       = errors.New("conflicting update statement")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotActivated       = errors.New("user not activated")
)
