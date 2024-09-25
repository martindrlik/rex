package relation_test

import (
	"fmt"
	"maps"
	"testing"

	"github.com/martindrlik/rex/relation"
)

func TestEqual(t *testing.T) {
	bu := func(count int) *relation.Relation {
		r := relation.New(build().userSchema...)
		for i := 1; i <= count; i++ {
			r.Add(build().someUser(i))
		}
		return r
	}

	u1 := bu(2)
	u2 := bu(2)

	if !relation.Equal(u1, u2) {
		t.Errorf("expected relations %+v and %+v to be equal", u1, u2)
	}
}

func TestClone(t *testing.T) {
	t.Run("IncludeAll", func(t *testing.T) {
		users := relation.New("name", "bornYear")
		users.Add(build().someUser(1))
		users.Add(build().someUser(2))
		clone := users.Clone(relation.IncludeAll)
		for _, user := range []map[string]any{build().someUser(1), build().someUser(2)} {
			if !clone.Has(user) {
				t.Errorf("expected clone to have %+v", user)
			}
		}
	})
	t.Run("IncludeOnly", func(t *testing.T) {
		users := relation.New("name", "bornYear")
		users.Add(build().someUser(1))
		users.Add(build().someUser(2))
		clone := users.Clone(relation.IncludeOnly(build().someUser(2)))
		if clone.Has(build().someUser(1)) {
			t.Errorf("expected clone to not have %+v", build().someUser(1))
		}
		if !clone.Has(build().someUser(2)) {
			t.Errorf("expected clone to have %+v", build().someUser(2))
		}
	})
}

func TestHas(t *testing.T) {
	users := relation.New("name", "bornYear")
	users.Add(build().someUser(1))
	if !users.Has(build().someUser(1)) {
		t.Errorf("expected users relation to have tuple %+v", build().someUser(1))
	}
	if users.Has(build().someUser(2)) {
		t.Errorf("expected users relation not to have tuple %+v", build().someUser(2))
	}
}

func build() struct {
	userSchema []string
	user       func(name string, bornYear int) map[string]any
	someUser   func(i int) map[string]any
} {
	user := func(name string, bornYear int) map[string]any {
		return map[string]any{"name": name, "bornYear": bornYear}
	}
	return struct {
		userSchema []string
		user       func(name string, bornYear int) map[string]any
		someUser   func(i int) map[string]any
	}{
		userSchema: []string{"name", "bornYear"},
		user:       user,
		someUser:   func(i int) map[string]any { return user(fmt.Sprintf("SomeUser%d", i), 2000+i) },
	}
}

func TestBuilder(t *testing.T) {
	actual := build().user("Foo", 2004)
	expect := map[string]any{"name": "Foo", "bornYear": int(2004)}
	if !maps.Equal(actual, expect) {
		t.Errorf("expected builder to create user %+v, got %+v", expect, actual)
	}
}

func must[V any](v V, err error) V {
	if err != nil {
		panic(err)
	}

	return v
}

func must1(err error) {
	if err != nil {
		panic(err)
	}
}
