package table

func (t *Table) Union(u *Table) *Table {
	x := New()
	for _, tuple := range t.tuples {
		x.tuples = append(x.tuples, tuple)
	}
	for _, tuple := range u.tuples {
		x.Add(tuple)
	}
	return x
}
