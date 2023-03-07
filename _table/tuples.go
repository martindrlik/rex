package table

type Tuples []map[string]any

func (tuples Tuples) Contains(tuple map[string]any) bool {
	for _, t := range tuples {
		if T(t).Equals(tuple) {
			return true
		}
	}
	return false
}
