package errors

import "errors"

var (
	ErrAuthAccessDenied = errors.New("access denied")

	ErrAuthGitHubUserNotReceived = errors.New("cannot receive user from github api")
	ErrAuthProviderNotFound      = errors.New("provider query parameter is missing")
)
