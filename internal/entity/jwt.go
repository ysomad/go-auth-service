package entity

import "github.com/google/uuid"

// JWT represents HTTP response data set of JSON web token authentication flow
type JWT struct {
	AccessToken  string    `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
	RefreshToken uuid.UUID `json:"refreshToken" example:"c84f18a2-c6c7-4850-be15-93f9cbaef3b3"`
	ExpiresIn    int       `json:"-"`
}
