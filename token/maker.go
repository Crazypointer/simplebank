package token

import "time"

type Maker interface {
	// CreateToken creates a new token for a specific username and with a specific duration.
	CreateToken(username string, duration time.Duration) (string, error)
	// VerifyToken checks if the token is valid or not.
	VerifyToken(token string) (*Payload, error)
}
