package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidURLFormat   = errors.New("invalid URL format")
	ErrURLFetchFailed     = errors.New("failed to fetch URL")
	ErrHTMLParseFailed    = errors.New("failed to parse HTML")
	ErrTokenInvalid       = errors.New("invalid or expired token")
	ErrMissingAuthHeader  = errors.New("missing or invalid Authorization header")
)
