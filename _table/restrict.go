package table

func (t *Table) Restrict(f func(map[string]any) bool) *Table {
	x := New()
	for _, tuple := range t.tuples {
		if f(tuple) {
			x.Add(tuple)
		}
	}
	return x
}
