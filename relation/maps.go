package relation

func makeMapFromSlice[K comparable](s ...K) map[K]struct{} {
	m := map[K]struct{}{}
	for _, k := range s {
		m[k] = struct{}{}
	}
	return m
}

func alwaysEqual[V1, V2 comparable](V1, V2) bool { return true }
