package relation

import (
	"iter"
	"maps"
)

type Relation struct {
	schema map[string]struct{}
	tuples []map[string]any
}

func New(schema ...string) *Relation {
	return &Relation{makeMapFromSlice(schema...), nil}
}

func Equal(r, s *Relation) bool {
	if !maps.Equal(r.schema, s.schema) {
		return false
	}
	if len(r.tuples) != len(s.tuples) {
		return false
	}
	for _, tup := range r.tuples {
		if !s.Has(tup) {
			return false
		}
	}
	return true
}

func (r *Relation) Clone(include func(map[string]any) bool) *Relation {
	s := &Relation{schema: r.schema}
	for tuple := range r.Tuples(include) {
		s.tuples = append(s.tuples, tuple)
	}
	return s
}

func (r *Relation) Add(tuple map[string]any) error {
	if !maps.EqualFunc(r.schema, tuple, alwaysEqual) {
		return ErrMismatch
	}
	if !r.Has(tuple) {
		r.tuples = append(r.tuples, tuple)
	}
	return nil
}

func (r *Relation) Tuples(include func(map[string]any) bool) iter.Seq[map[string]any] {
	return func(yield func(map[string]any) bool) {
		for _, tup := range r.tuples {
			if include(tup) && !yield(tup) {
				return
			}
		}
	}
}

func IncludeAll(map[string]any) bool { return true }

func IncludeOnly(tuple map[string]any) func(map[string]any) bool {
	return func(tup map[string]any) bool {
		return maps.Equal(tuple, tup)
	}
}

func (r *Relation) Has(tuple map[string]any) bool {
	for range r.Tuples(IncludeOnly(tuple)) {
		return true
	}
	return false
}
