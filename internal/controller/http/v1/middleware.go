package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/internal/service"
)

func sessionMiddleware(s service.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie(sessionIDKey)
		if err != nil {
			abortWithError(c, http.StatusUnauthorized, entity.ErrUnauthorized)
			return
		}

		// TODO: add session id validation

		ctx := c.Request.Context()

		sess, err := s.Find(ctx, sid)
		if err != nil {
			abortWithError(c, http.StatusUnauthorized, entity.ErrSessionExpired)
			return
		}

		d, err := entity.NewDevice(c.Request.Header.Get("User-Agent"), c.ClientIP())
		if err != nil {
			abortWithError(c, http.StatusUnauthorized, entity.ErrSessionExpired)
			return
		}

		// Check current request device vs session device
		if sess.UserIP != d.UserIP || sess.UserAgent != d.UserAgent {
			// TODO: send notification that someone logged in on new device
			s.Terminate(ctx, sid)

			abortWithError(c, http.StatusUnauthorized, entity.ErrSessionExpired)
			return
		}

		c.Next()
	}
}

/*
func getUserID(c *gin.Context) (uuid.UUID, error) {
	idStr := c.GetString("user")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func jwtMiddleware(jwt auth.JWT) gin.HandlerFunc {
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
		userID, err := jwt.Validate(accessToken)
		if err != nil {
			abortWithError(c, http.StatusUnauthorized, err)
			return
		}

		c.Set("user", userID)
		c.Next()
	}
}


*/
