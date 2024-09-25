package table

func (t *Table) Project(schema ...string) *Table {
	if !t.isSchemaSubset(schema) {
		return New()
	}
	x := New()
	add := func() func(map[string]any) {
		ptf := projectTuple(schema...)
		return func(tuple map[string]any) {
			tuple, ok := ptf(tuple)
			if ok {
				x.Add(tuple)
			}
		}
	}()
	for _, tuple := range t.tuples {
		add(tuple)
	}
	return x
}

func (t *Table) isSchemaSubset(schema []string) bool {
	x := t.Schema()
	if len(x) < len(schema) {
		return false
	}
	for _, k := range schema {
		if _, ok := x[k]; !ok {
			return false
		}
	}
	return true
}

func projectTuple(schema ...string) func(map[string]any) (map[string]any, bool) {
	return func(tuple map[string]any) (map[string]any, bool) {
		x := make(map[string]any)
		for _, k := range schema {
			v, ok := tuple[k]
			if ok {
				x[k] = v
			}
		}
		return x, len(x) > 0
	}
}
