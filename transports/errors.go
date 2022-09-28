package transports

import "errors"

// ErrNoClientSet is when no client is set
var ErrNoClientSet = errors.New("no transport client set")
var ErrFailedLogin = errors.New("failed to login to server")
