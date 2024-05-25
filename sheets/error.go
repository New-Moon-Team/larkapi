package sheets

import (
	"errors"
)

var (
	ErrSheetNotFound       = errors.New("Sheet not found")
	ErrNotAuthenticated    = errors.New("Not authenticated")
	ErrLoginFailed         = errors.New("Login failed")
	ErrUnknownAuthResponse = errors.New("Unknown auth response")
)
