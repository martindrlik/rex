package relation

import "errors"

var (
	ErrAlreadyPresent = errors.New("already present")
	ErrMissingSchema  = errors.New("missing schema")
	ErrSchemaMismatch = errors.New("schema mismatch")
)
