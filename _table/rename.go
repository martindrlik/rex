package table

func (t *Table) Rename(rename map[string]string) *Table {
	x := New()
	for _, tuple := range t.tuples {
		y := make(map[string]any)
		for k, v := range tuple {
			if k, ok := rename[k]; ok {
				y[k] = v
				continue
			}
			y[k] = v
		}
		x.Add(y)
	}
	return x
}
