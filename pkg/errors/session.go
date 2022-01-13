package errors

import "errors"

var (
	ErrSessionNotFound         = errors.New("session not found")
	ErrSessionExpired          = errors.New("session expired")
	ErrSessionNotCreated       = errors.New("error occured during session creation")
	ErrSessionNotTerminated    = errors.New("current session cannot be terminated, use logout instead")
	ErrSessionMismatchedDevice = errors.New("session does not match with device")
	ErrSessionContextNotFound  = errors.New("session not found in context")
)
