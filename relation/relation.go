package relation

import (
	"iter"
	"maps"
	"slices"
)

type Relation struct {
	schema map[string]struct{}
	tuples []map[string]any
}

func New(attributes ...string) (*Relation, error) {
	if len(attributes) == 0 {
		return nil, ErrMissingSchema
	}

	schema := make(map[string]struct{}, len(attributes))
	for _, attribute := range attributes {
		schema[attribute] = struct{}{}
	}

	tuples := make([]map[string]any, 0)
	return &Relation{
		schema,
		tuples,
	}, nil
}

func (r *Relation) Add(tuple map[string]any) error {
	if !r.hasMatchingSchema(tuple) {
		return ErrSchemaMismatch
	}

	if r.Contain(tuple) {
		return ErrAlreadyPresent
	}

	r.tuples = append(r.tuples, tuple)
	return nil
}

func (r *Relation) Contain(tuple map[string]any) bool {
	if r.hasMatchingSchema(tuple) {
		for t1 := range r.List() {
			if maps.Equal(t1, tuple) {
				return true
			}
		}
	}
	return false
}

func (r *Relation) List() iter.Seq[map[string]any] {
	return func(yield func(map[string]any) bool) {
		for _, t := range r.tuples {
			if !yield(t) {
				return
			}
		}
	}
}

func Union(r1, r2 *Relation) (*Relation, error) {
	if !maps.Equal(r1.schema, r2.schema) {
		return nil, ErrSchemaMismatch
	}
	r := &Relation{
		schema: map[string]struct{}{},
		tuples: make([]map[string]any, len(r1.tuples)+len(r2.tuples))}
	maps.Copy(r.schema, r1.schema)
	r.tuples = slices.Clone(r1.tuples)
	for t1 := range r2.List() {
		if !r.Contain(t1) {
			r.tuples = append(r.tuples, t1)
		}
	}
	return r, nil
}

func (r *Relation) hasMatchingSchema(tuple map[string]any) bool {
	if len(r.schema) == len(tuple) {
		for attribute := range tuple {
			if _, ok := r.schema[attribute]; !ok {
				return false
			}
		}

		return true
	}

	return false
}
