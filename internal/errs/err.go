package errs

import "errors"

var (
	ErrDuplicateUrl   = errors.New("duplicate url")
	ErrUrlNotFound    = errors.New("url not found")
	ErrInternalServer = errors.New("internal server error")
)
