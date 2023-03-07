package table

import (
	"errors"
	"fmt"
)

var (
	ErrAlreadyPresent     = errors.New("already present")
	ErrDuplicateAttribute = func(attribute string) error { return fmt.Errorf("duplicate attribute: %s", attribute) }
	ErrMissingSchema      = errors.New("missing schema")
	ErrSchemaMismatch     = errors.New("schema mismatch")
)
