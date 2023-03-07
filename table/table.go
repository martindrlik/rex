package table

import (
	"fmt"
	"iter"
	"maps"
)

type Table struct {
	schema    map[string]struct{}
	relations []*Relation
}

func New(aa ...string) (*Table, error) {
	if len(aa) == 0 {
		return nil, ErrMissingSchema
	}

	schema := make(map[string]struct{}, len(aa))
	for _, a := range aa {
		schema[a] = struct{}{}
	}

	return &Table{schema, []*Relation{}}, nil
}

// Count returns number of tuples in the table.
func (b1 *Table) Count() int {
	n := 0
	for _, r1 := range b1.relations {
		n += r1.Count()
	}
	return n
}

func (b1 *Table) Add(t1 tup) error {
	if !b1.hasMatchingSchema(t1) {
		return ErrSchemaMismatch
	}

	for _, r1 := range b1.relations {
		err := r1.Add(t1)
		if err == ErrSchemaMismatch {
			continue
		}
		return err
	}

	aa := make([]string, 0, len(t1))
	for a := range t1 {
		aa = append(aa, a)
	}
	r1, err := NewRelation(aa...)
	if err != nil {
		return err
	}
	b1.relations = append(b1.relations, r1)
	return b1.Add(shallowCopyTuple(t1))
}

func (b1 *Table) Contain(t1 tup) bool {
	for _, r1 := range b1.relations {
		if r1.Contain(t1) {
			return true
		}
	}
	return false
}

// List returns a sequence of tuples in table.
func (b1 *Table) List() iter.Seq[tup] {
	return func(yield func(tup) bool) {
		for _, r1 := range b1.relations {
			for t1 := range r1.List() {
				if !yield(t1) {
					return
				}
			}
		}
	}
}

// Union creates new table with all tuples from both involved tables.
func (b1 *Table) Union(b2 *Table) (*Table, error) {
	if !maps.Equal(b1.schema, b2.schema) {
		return nil, ErrSchemaMismatch
	}

	b3 := newTable(b1.schema, len(b1.relations)+len(b2.relations))
	for _, r1 := range b1.relations {
		for t1 := range r1.List() {
			if err := b3.Add(t1); err != nil {
				panic(fmt.Errorf("unexpected error %w", err))
			}
		}
	}
	for _, r2 := range b2.relations {
		for t2 := range r2.List() {
			if !b3.Contain(t2) {
				if err := b3.Add(t2); err != nil {
					panic(fmt.Errorf("unexpected error %w", err))
				}
			}
		}
	}

	return b3, nil
}

func (b1 *Table) hasMatchingSchema(t1 tup) bool {
	if len(t1) > len(b1.schema) {
		return false
	}

	for a := range t1 {
		if _, ok := b1.schema[a]; !ok {
			return false
		}
	}

	return true
}

func newTable(srcSchema map[string]struct{}, relationCapacity int) *Table {
	dstSchema := make(map[string]struct{}, len(srcSchema))
	maps.Copy(dstSchema, srcSchema)
	return &Table{dstSchema, make([]*Relation, 0, relationCapacity)}
}
