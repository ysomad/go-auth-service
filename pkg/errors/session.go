package errors

import "errors"

const (
	SessionNotFound   = "session not found"
	SessionExpired    = "session expired"
	SessionNotCreated = "error occured during session creation"
)

var (
	ErrSessionNotFound   = errors.New(SessionNotFound)
	ErrSessionExpired    = errors.New(SessionExpired)
	ErrSessionNotCreated = errors.New(SessionNotCreated)
)
