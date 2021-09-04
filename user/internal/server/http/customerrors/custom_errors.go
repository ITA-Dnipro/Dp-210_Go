package customerrors

import "errors"

var (
	NotFound = errors.New("not found")
	Internal = errors.New("internal")
)
