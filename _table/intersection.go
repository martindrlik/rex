package table

func (t *Table) Intersection(u *Table) *Table {
	x := t.Difference(u)
	return u.Difference(x)
}
