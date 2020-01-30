package server

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) auth(c *gin.Context) {
	username, password, err := parseBasicAuth(c)
	if err != nil {
		abortUnauthorized(c, err)
		return
	}

	id := c.Param("id")
	if id != username {
		abortUnauthorized(c, errors.New("invalid username"))
		return
	}

	app, err := s.store.GetApp(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if ok := app.VerifySecret(password); !ok {
		abortUnauthorized(c, errors.New("credential mismatch"))
		return
	}
}

func abortUnauthorized(c *gin.Context, err error) {
	c.Header("WWW-Authenticate", `Basic realm="Authorization Required"`)
	c.AbortWithError(http.StatusUnauthorized, err)
}

func parseBasicAuth(c *gin.Context) (string, string, error) {
	header := c.GetHeader("Authorization")

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid authorization header")
	}
	if method := strings.ToLower(parts[0]); method != "basic" {
		return "", "", errors.New("not basic auth")
	}

	cred, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", errors.New("invalid credential encoding")
	}

	parts = strings.SplitN(string(cred), ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid credential format")
	}

	return parts[0], parts[1], nil
}
