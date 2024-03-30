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

	ErrBucketNotFound      = errors.New("bucket not found")
	ErrBucketAlreadyExists = errors.New("bucket already exists")
	ErrBucketConvertError  = errors.New("ups, not bucket value")

	ErrBucketRateInvalid = errors.New("rate can't be zero value")
)
