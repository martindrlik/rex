package table

func (t *Table) Difference(u *Table) *Table {
	x := New()
	isComplete := t.isCompleteTuple()
	for _, tuple := range t.tuples {
		if !isComplete(tuple) || !u.CompleteTuples().Contains(tuple) {
			x.Add(tuple)
		}
	}
	return x
}
