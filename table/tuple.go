package table

import "maps"

type tup = map[string]any

func shallowCopyTuple(t1 tup) tup {
	t2 := make(tup, len(t1))
	maps.Copy(t2, t1)
	return t2
}
