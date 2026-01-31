package errors

import (
	"errors"
)

var (
	ErrTaskNotFound   = errors.New("task not found")
	ErrInternalServer = errors.New("Internal server error")
)
