package apperrors

import "errors"

var (
	ErrCSRFTokenPublicNotFound = errors.New("csrf token not found in request headers or query")
	ErrCSRFTokenCookieNotFound = errors.New("csrf token not found in cookies")
	ErrCSRFDetected            = errors.New("csrf tokens in headers and cookies are not the same")
)
