package table

import (
	"iter"
	"maps"
)

// Relation represents a relation with a schema and tuples.
type Relation struct {
	schema map[string]struct{}
	tuples []tup
}

// NewRelation creates a new relation with attributes given by aa.
func NewRelation(aa ...string) (*Relation, error) {
	if len(aa) == 0 {
		return nil, ErrMissingSchema
	}

	schema := make(map[string]struct{}, len(aa))
	for _, a := range aa {
		schema[a] = struct{}{}
	}

	return &Relation{schema, []tup{}}, nil
}

// Count returns number of tuples in the relation.
func (r1 *Relation) Count() int { return len(r1.tuples) }

// Add adds tuple t1 to relation. Returns error if tuple t1 has different schema or if already present.
func (r1 *Relation) Add(t1 map[string]any) error {
	if !r1.hasMatchingSchema(t1) {
		return ErrSchemaMismatch
	}

	if r1.Contain(t1) {
		return ErrAlreadyPresent
	}

	r1.tuples = append(r1.tuples, t1)
	return nil
}

// Contain returns true if tuple t2 is present in relation. Returns false otherwise.
func (r1 *Relation) Contain(t2 map[string]any) bool {
	if r1.hasMatchingSchema(t2) {
		for t1 := range r1.List() {
			if maps.Equal(t1, t2) {
				return true
			}
		}
	}
	return false
}

// List returns a sequence of tuples in relation.
func (r1 *Relation) List() iter.Seq[map[string]any] {
	return func(yield func(tup) bool) {
		for _, t1 := range r1.tuples {
			if !yield(t1) {
				return
			}
		}
	}
}

// Union creates new relation with all tuples from both involved relations.
func (r1 *Relation) Union(r2 *Relation) (*Relation, error) {
	if !maps.Equal(r1.schema, r2.schema) {
		return nil, ErrSchemaMismatch
	}

	r3 := newRelation(r1.schema, len(r1.tuples)+len(r2.tuples))

	for t1 := range r1.List() {
		r3.tuples = append(r3.tuples, shallowCopyTuple(t1))
	}

	for t2 := range r2.List() {
		if !r3.Contain(t2) {
			r3.tuples = append(r3.tuples, shallowCopyTuple(t2))
		}
	}
	return r3, nil
}

func (r1 *Relation) Difference(r2 *Relation) (*Relation, error) {
	if !maps.Equal(r1.schema, r2.schema) {
		return nil, ErrSchemaMismatch
	}

	r3 := newRelation(r1.schema, len(r1.tuples))

	for t1 := range r1.List() {
		if !r2.Contain(t1) {
			r3.tuples = append(r3.tuples, shallowCopyTuple(t1))
		}
	}

	return r3, nil
}

func (r1 *Relation) NaturalJoin(r2 *Relation) *Relation {
	common := r1.common(r2)
	concat := func(t, t2 tup) (tup, bool) {
		t3 := make(tup, len(t)+len(t2))
		for a := range common {
			if t[a] != t2[a] {
				return nil, false
			}
		}
		maps.Copy(t3, t)
		maps.Copy(t3, t2)
		return t3, true
	}

	r3 := &Relation{schema: map[string]struct{}{}}
	maps.Copy(r3.schema, r1.schema)
	maps.Copy(r3.schema, r2.schema)

	for t1 := range r1.List() {
		for t2 := range r2.List() {
			if t3, ok := concat(t1, t2); ok {
				r3.Add(t3)
			}
		}
	}

	return r3
}

func (r1 *Relation) Rename(f func(old string) (new string)) (*Relation, error) {
	schema := make(map[string]struct{}, len(r1.schema))
	schemaCount := make(map[string]int, len(r1.schema))

	for a1 := range r1.schema {
		a := f(a1)
		schemaCount[a]++
		if schemaCount[a] > 1 {
			return nil, ErrDuplicateAttribute(a)
		}
		schema[a] = struct{}{}
	}

	r2 := newRelation(schema, len(r1.tuples))

	for t1 := range r1.List() {
		t2 := make(tup, len(t1))
		for a1, v1 := range t1 {
			t2[f(a1)] = v1
		}
		r2.tuples = append(r2.tuples, t2)
	}

	return r2, nil
}

func (r1 *Relation) Intersection(r2 *Relation) (*Relation, error) {
	if !maps.Equal(r1.schema, r2.schema) {
		return nil, ErrSchemaMismatch
	}

	r3 := newRelation(r1.schema, 0)

	for t1 := range r1.List() {
		if r2.Contain(t1) {
			r3.tuples = append(r3.tuples, shallowCopyTuple(t1))
		}
	}

	return r3, nil
}

func (r1 *Relation) Project(aa ...string) (*Relation, error) {
	if len(aa) == 0 {
		return nil, ErrMissingSchema
	}

	schema := make(map[string]struct{}, len(aa))
	for _, a := range aa {
		if _, ok := r1.schema[a]; !ok {
			return nil, ErrSchemaMismatch
		}
		schema[a] = struct{}{}
	}

	r2 := newRelation(schema, len(r1.tuples))

	for t1 := range r1.List() {
		t2 := make(tup, len(schema))
		for a := range schema {
			t2[a] = t1[a]
		}
		r2.tuples = append(r2.tuples, t2)
	}

	return r2, nil
}

func (r1 *Relation) common(r2 *Relation) map[string]struct{} {
	m := make(map[string]struct{})
	for a := range r1.schema {
		if _, ok := r2.schema[a]; ok {
			m[a] = struct{}{}
		}
	}
	return m
}

func (r1 *Relation) hasMatchingSchema(t2 map[string]any) bool {
	if len(r1.schema) == len(t2) {
		for a := range t2 {
			if _, ok := r1.schema[a]; !ok {
				return false
			}
		}

		return true
	}

	return false
}

func newRelation(srcSchema map[string]struct{}, tupleCapacity int) *Relation {
	dstSchema := make(map[string]struct{}, len(srcSchema))
	maps.Copy(dstSchema, srcSchema)
	return &Relation{dstSchema, make([]tup, 0, tupleCapacity)}
}
