package relation

import (
	"iter"
	"maps"
)

type tuple = map[string]any

// Relation represents a relation with a schema and tuples.
type Relation struct {
	schema map[string]struct{}
	tuples []tuple
}

// New creates a new relation with given attributes.
func New(attributes ...string) (*Relation, error) {
	if len(attributes) == 0 {
		return nil, ErrMissingSchema
	}

	schema := make(map[string]struct{}, len(attributes))
	for _, attribute := range attributes {
		schema[attribute] = struct{}{}
	}

	tuples := make([]tuple, 0)
	return &Relation{
		schema,
		tuples,
	}, nil
}

// Count returns number of tuples in relation.
func (r1 *Relation) Count() int { return len(r1.tuples) }

// Add tuple t1 to relation. Returns error if tuple has different schema or if already present.
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
	return func(yield func(tuple) bool) {
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

	r3 := &Relation{schema: make(map[string]struct{}, len(r1.schema)), tuples: make([]tuple, 0, len(r1.tuples)+len(r2.tuples))}
	maps.Copy(r3.schema, r1.schema)

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

// Difference creates new relation with only tuples that are not included in r2.
func (r1 *Relation) Difference(r2 *Relation) (*Relation, error) {
	if !maps.Equal(r1.schema, r2.schema) {
		return nil, ErrSchemaMismatch
	}

	r3 := &Relation{schema: map[string]struct{}{}, tuples: make([]tuple, 0, len(r1.tuples))}

	for t1 := range r1.List() {
		if !r2.Contain(t1) {
			r3.tuples = append(r3.tuples, shallowCopyTuple(t1))
		}
	}

	return r3, nil
}

func (r1 *Relation) NaturalJoin(r2 *Relation) *Relation {
	common := r1.commonAttributes(r2)
	concat := func(t, t2 tuple) (tuple, bool) {
		t3 := make(tuple, len(t)+len(t2))
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

func (r1 *Relation) Rename(newByOld map[string]string) (*Relation, error) {

	// validate new schema
	schema := make(map[string]struct{}, len(r1.schema))
	schemaCount := make(map[string]int, len(r1.schema))
	for a1 := range r1.schema {
		var a string
		if a2, ok := newByOld[a1]; ok {
			a = a2
		} else {
			a = a1
		}
		schemaCount[a]++
		if schemaCount[a] > 1 {
			return nil, ErrDuplicateAttribute(a)
		}
		schema[a] = struct{}{}
	}

	r2 := &Relation{schema: map[string]struct{}{}, tuples: make([]tuple, 0, len(r1.tuples))}
	maps.Copy(r2.schema, schema)

	for t1 := range r1.List() {
		t2 := make(tuple, len(t1))
		for a1, v1 := range t1 {
			var a string
			if a2, ok := newByOld[a1]; ok {
				a = a2
			} else {
				a = a1
			}
			t2[a] = v1
		}
		r2.tuples = append(r2.tuples, t2)
	}

	return r2, nil
}

func (r1 *Relation) commonAttributes(r2 *Relation) map[string]struct{} {
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
		for attribute := range t2 {
			if _, ok := r1.schema[attribute]; !ok {
				return false
			}
		}

		return true
	}

	return false
}

func shallowCopyTuple(t1 tuple) tuple {
	t2 := make(tuple, len(t1))
	maps.Copy(t2, t1)
	return t2
}
