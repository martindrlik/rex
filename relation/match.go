package relation

import "maps"

func Tuple(tuple map[string]any) func(map[string]any) bool {
	return func(m map[string]any) bool {
		return maps.Equal(tuple, m)
	}
}
