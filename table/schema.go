package table

import "errors"

func validateSchema(schema map[string]struct{}) error {
	if len(schema) == 0 {
		return errors.New("missing some schema attributes")
	}
	return nil
}
