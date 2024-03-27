package errors

import (
	"errors"
)

var (
	ErrInvalidIP       = errors.New("invalid ip")
	ErrInvalidSubnet   = errors.New("invalid subnet")
	ErrInvalidLogin    = errors.New("invalid login")
	ErrInvalidPassword = errors.New("invalid password")

	ErrSubnetNotFound = errors.New("subnet or ip not found")
)
