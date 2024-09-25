package table

import "reflect"

func AreTuplesEqual(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k, av := range a {
		if bv, ok := b[k]; !ok {
			return false
		} else if !reflect.DeepEqual(av, bv) {
			return false
		}
	}
	return true
}
