package table

import (
	"golang.org/x/exp/maps"
)

func (t *Table) NaturalJoin(u *Table) *Table {
	equal := func() func(a, b map[string]any) bool {
		common := map[string]struct{}{}
		ts, us := t.Schema(), u.Schema()
		for k := range ts {
			if _, ok := us[k]; ok {
				common[k] = struct{}{}
			}
		}
		if len(common) == 0 {
			return func(a, b map[string]any) bool { return true }
		}
		return func(a, b map[string]any) bool {
			return T(a).equalsOn(b, maps.Keys(common)...)
		}
	}()
	x := New()
	for _, tuple := range t.CompleteTuples() {
		for _, other := range u.CompleteTuples() {
			if equal(tuple, other) {
				merged := T(tuple).Merge(other)
				x.Add(merged)
			}
		}
	}
	return x
}

func (t *Table) hasCommon(u *Table) bool {
	for k := range t.Schema() {
		if _, ok := u.Schema()[k]; ok {
			return true
		}
	}
	return false
}
