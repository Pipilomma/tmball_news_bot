package errs

import "errors"

var (
	ErrUserUlreadySub = errors.New("user is ulready sub")
	ErrNewsNotFound   = errors.New("zero news with that title")
)
