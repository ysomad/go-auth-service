package entity

import (
	"github.com/google/uuid"
)

type (
	LoginRequest struct {
		Email       string `json:"email" example:"user@mail.com" binding:"required,email,lte=255"`
		Password    string `json:"password" example:"secret" binding:"required,gte=6,lte=128"`
		Fingerprint string `json:"fingerprint" example:"c84f18a2-c6c7-4850-be15-93f9cbaef3b3" binding:"required,uuid4"`
	}

	LoginResponse struct {
		AccessToken  string    `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
		RefreshToken uuid.UUID `json:"refreshToken" example:"c84f18a2-c6c7-4850-be15-93f9cbaef3b3"`
		ExpiresIn    int       `json:"-"`
	}
)
