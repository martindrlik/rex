package table

import "reflect"

type (
	T map[string]any
)

func (t T) Equals(u T) bool {
	return len(t) == len(u) && t.equals(u)
}

func (t T) equalsOn(u T, attributes ...string) bool {
	for _, attribute := range attributes {
		tv, tok := t[attribute]
		uv, uok := u[attribute]
		if !tok || !uok {
			return false
		}
		return reflect.DeepEqual(tv, uv)
	}
	return true
}

func (t T) equals(u T) bool {
	for k, v := range t {
		w, ok := u[k]
		if !ok || !reflect.DeepEqual(v, w) {
			return false
		}
	}
	return true
}

func (t T) Merge(u T) T {
	x := make(T)
	for k, v := range t {
		x[k] = v
	}
	for k, v := range u {
		x[k] = v
	}
	return x
}
