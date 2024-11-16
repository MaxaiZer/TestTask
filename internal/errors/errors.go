package errors

import (
	"github.com/pkg/errors"
)

var UnsupportedOperation = errors.New("operation unsupported")
var InsufficientBalance = errors.New("insufficient balance")
var TooManyRequests = errors.New("too many requests")
var NotFound = errors.New("not found")
