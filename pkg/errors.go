package pkg

import (
	"errors"
)

var (
	ErrSomething        = errors.New("something went wrong")
	ErrMethodNotAllowed = errors.New("method not allowed")
)
