package apperrors

import "errors"

var (
	ErrAuthAccessDenied          = errors.New("access denied")
	ErrAuthProviderNotFound      = errors.New("provider query parameter is missing")
	ErrAuthGitHubUserNotReceived = errors.New("cannot receive user from github api")
)
