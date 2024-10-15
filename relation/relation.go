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
		for t := range r.List() {
			if maps.Equal(t, tuple) {
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
		tuples: slices.Clone(r1.tuples)}
	maps.Copy(r.schema, r1.schema)
	for t2 := range r2.List() {
		if !r.Contain(t2) {
			r.tuples = append(r.tuples, t2)
		}
	}
	return r, nil
}

// Difference creates new relation with only tuples that are not included in r2.
func (r *Relation) Difference(r2 *Relation) (*Relation, error) {
	if !maps.Equal(r.schema, r2.schema) {
		return nil, ErrSchemaMismatch
	}

	r3 := &Relation{
		schema: map[string]struct{}{},
		tuples: make([]map[string]any, 0, len(r.tuples))}

	for t := range r.List() {
		if !r2.Contain(t) {
			t3 := make(map[string]any, len(t))
			maps.Copy(t3, t)
			r3.tuples = append(r3.tuples, t3)
		}
	}

	return r3, nil
}

// Count return number of tuples in relation.
func (r *Relation) Count() int { return len(r.tuples) }

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
