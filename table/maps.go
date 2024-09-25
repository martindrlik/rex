package table

func hasSameKeys[K comparable, U, V any](a map[K]U, b map[K]V) bool {
	if len(a) != len(b) {
		return true
	}
	for k := range a {
		if _, ok := b[k]; !ok {
			return false
		}
	}
	return true
}
