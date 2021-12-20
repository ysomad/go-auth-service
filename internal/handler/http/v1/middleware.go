package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/service"
)

func sessionMiddleware(s service.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie("id") // TODO: refactor
		if err != nil {
			// TODO: return generic err pkg httperror
			abortWithError(c, http.StatusUnauthorized, domain.ErrUnauthorized)
			return
		}

		// TODO: add session id validation

		ctx := c.Request.Context()

		sess, err := s.Get(ctx, sid)
		if err != nil {
			// TODO: return generic err pkg httperror
			abortWithError(c, http.StatusUnauthorized, domain.ErrSessionExpired)
			return
		}

		d, err := domain.NewDevice(c.Request.Header.Get("User-Agent"), c.ClientIP())
		if err != nil {
			// TODO: return generic err pkg httperror
			abortWithError(c, http.StatusUnauthorized, domain.ErrSessionExpired)
			return
		}

		// Check current request device vs session device
		if sess.IP != d.IP || sess.UserAgent != d.UserAgent {
			// TODO: send notification that someone logged in on new device
			s.Terminate(ctx, sid)

			// TODO: return generic err pkg httperror
			abortWithError(c, http.StatusUnauthorized, domain.ErrSessionExpired)
			return
		}

		c.Next()
	}
}

func accountID(c *gin.Context) (string, error) {
	aid := c.GetString("user") // TODO: refactor

	_, err := uuid.Parse(aid)
	if err != nil {
		// TODO: return generic err pkg httperror
		return "", err
	}

	return aid, nil
}

// TODO: implement token middleware
/*
func tokenMiddleware(jwt auth.JWT) gin.HandlerFunc {
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
