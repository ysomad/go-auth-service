package v1

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/pkg/auth"
)

func authMiddleware(jwt auth.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if len(header) == 0 {
			abortWithError(c, http.StatusUnauthorized, errors.New("authorization header is empty"))
			return
		}

		fields := strings.Fields(header)
		if len(fields) < 2 {
			abortWithError(c, http.StatusUnauthorized, errors.New("invalid authorization header format"))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != "bearer" {
			abortWithError(c, http.StatusUnauthorized, errors.New("unsupported authorization type"))
			return
		}

		accessToken := fields[1]
		payload, err := jwt.Validate(accessToken)
		if err != nil {
			abortWithError(c, http.StatusUnauthorized, err)
			return
		}

		c.Set("user", payload["sub"].(string))
		c.Next()
	}
}

func getUserID(c *gin.Context) (uuid.UUID, error) {
	sub, ok := c.Get("user")
	if !ok {
		return uuid.UUID{}, errors.New("empty subject")
	}

	idString, ok := sub.(string)
	if !ok {
		return uuid.UUID{}, errors.New("invalid type of subject")
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}
