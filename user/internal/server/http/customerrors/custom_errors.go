package customerrors

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrInternal    = errors.New("internal")
	ErrDublication = errors.New("dublication one of the key")
	ErrForeignKey  = errors.New("foreign key constraint")
)
