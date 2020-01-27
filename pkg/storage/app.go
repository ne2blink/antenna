package storage

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// App represnets an application
type App struct {
	ID     string
	Name   string
	Secret string
}

// SetSecret generates a random secret.
func (a *App) SetSecret(secret string) (string, error) {
	if secret == "" {
		r := make([]byte, 16)
		rand.Read(r)
		secret = base64.RawURLEncoding.EncodeToString(r)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return secret, err
	}
	a.Secret = string(password)

	return secret, nil
}

// VerifySecret verifies secret.
func (a *App) VerifySecret(secret string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.Secret), []byte(secret)) == nil
}
