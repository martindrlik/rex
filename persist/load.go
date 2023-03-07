package persist

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/martindrlik/rex/table"
)

func Load(r io.Reader) (*table.Table, error) {
	dec := json.NewDecoder(r)
	tuples := []map[string]any{}
	if err := dec.Decode(&tuples); err != nil {
		return nil, err
	}
	return table.New().Add(tuples...), nil
}

// LoadSchemaMode loads table with schema and tuples.
// {"schema": [...], "tuples": [{...}, ...]}
func LoadSchemaMode(r io.Reader) (*table.Table, error) {
	dec := json.NewDecoder(r)
	raw := struct {
		Schema []string
		Tuples []map[string]any
	}{}
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}
	schema, m := func() ([]string, map[string]struct{}) {
		x := []string{}
		y := map[string]struct{}{}
		for _, attribute := range raw.Schema {
			if _, ok := y[attribute]; !ok {
				x = append(x, attribute)
				y[attribute] = struct{}{}
			}
		}
		return x, y
	}()
	isValidTuple := func(tuple map[string]any) bool {
		for k := range tuple {
			if _, ok := m[k]; !ok {
				return false
			}
		}
		return true
	}
	if len(schema) == 0 {
		return nil, errors.New("invalid schema")
	}
	tuples := func() []map[string]any {
		tuples := []map[string]any{}
		for _, tuple := range raw.Tuples {
			if isValidTuple(tuple) {
				tuples = append(tuples, tuple)
			}
		}
		return tuples
	}()
	return table.New(schema...).Add(tuples...), nil
}
