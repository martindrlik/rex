package relation_test

import (
	"maps"
	"testing"

	"github.com/martindrlik/rex/relation"
)

func TestNew(t *testing.T) {
	t.Run("new", testNew)
	t.Run("missing schema", testMissingSchema)
}

func testNew(t *testing.T) {
	_, err := relation.New("foo")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func testMissingSchema(t *testing.T) {
	_, err := relation.New()
	if err != relation.ErrMissingSchema {
		t.Errorf("(1) expected error %v, got %v", relation.ErrMissingSchema, err)
	}
	_, err = relation.New([]string{}...)
	if err != relation.ErrMissingSchema {
		t.Errorf("(2) expected error %v, got %v", relation.ErrMissingSchema, err)
	}
}

func TestAdd(t *testing.T) {
	t.Run("add", testAdd)
	t.Run("add multiple", testAddMultiple)
	t.Run("add duplicate error", testAddDuplicate)
	t.Run("not matching schema", testAddNotMatchingSchema)
}

func testAdd(t *testing.T) {
	r := newRelation("foo", "bar", "baz")
	if err := r.Add(map[string]any{"foo": 1, "bar": 2, "baz": 3}); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func testAddMultiple(t *testing.T) {
	r := newRelation("foo")
	if err := r.Add(map[string]any{"foo": 1}); err != nil {
		t.Errorf("(1) expected no error, got %v", err)
	}
	if err := r.Add(map[string]any{"foo": 2}); err != nil {
		t.Errorf("(2) expected no error, got %v", err)
	}
}

func testAddDuplicate(t *testing.T) {
	r := newRelation("foo")
	if err := r.Add(map[string]any{"foo": 1}); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if err := r.Add(map[string]any{"foo": 1}); err != relation.ErrAlreadyPresent {
		t.Errorf("expected error %v, got %v", relation.ErrAlreadyPresent, err)
	}
}

func testAddNotMatchingSchema(t *testing.T) {
	r := newRelation("foo")
	t1 := map[string]any{"bar": "baz"}
	t2 := map[string]any{"foo": "1", "bar": "2"}
	t3 := map[string]any{"foo": "bar"}
	if err := r.Add(t1); err != relation.ErrSchemaMismatch {
		t.Errorf("(1) expected error %v, got %v", relation.ErrSchemaMismatch, err)
	}
	if err := r.Add(t2); err != relation.ErrSchemaMismatch {
		t.Errorf("(2) expected error %v, got %v", relation.ErrSchemaMismatch, err)
	}
	if err := r.Add(t3); err != nil {
		t.Errorf("(3) expected to have matching schema, got error %v", err)
	}
}

func TestContain(t *testing.T) {
	r := newRelation("foo", "bar")
	if err := r.Add(map[string]any{"foo": 1, "bar": 2}); err != nil {
		panic(err)
	}
	t.Run("true", func(t *testing.T) {
		if !r.Contain(map[string]any{"foo": 1, "bar": 2}) {
			t.Errorf("expected relation to contain %+v", map[string]any{"foo": 1, "bar": 2})
		}
	})
	t.Run("false different schema", func(t *testing.T) {
		if r.Contain(map[string]any{"foo": 1}) {
			t.Error("expected relation to not contain tuple due to different schema")
		}
	})
	t.Run("false different value", func(t *testing.T) {
		if r.Contain(map[string]any{"foo": 1, "bar": 3}) {
			t.Error("expected relation to not contain tuple due to different value")
		}
	})
	t.Run("false different values", func(t *testing.T) {
		if r.Contain(map[string]any{"foo": 2, "bar": 3}) {
			t.Error("expected relation to not contain tuple due to different values")
		}
	})
}

func TestList(t *testing.T) {
	t.Run("List", testList)
}

func testList(t *testing.T) {
	r := newRelation("foo", "bar")
	if err := r.Add(map[string]any{"foo": 1, "bar": 2}); err != nil {
		panic(err)
	}
	if err := r.Add(map[string]any{"foo": 3, "bar": 4}); err != nil {
		panic(err)
	}
	expect := []map[string]any{
		{"foo": 1, "bar": 2},
		{"foo": 3, "bar": 4},
	}
	idx := 0
	for t1 := range r.List() {
		if !maps.Equal(t1, expect[idx]) {
			t.Errorf("(%d) expected tuples %+v and %+v to be equal", idx, t1, expect[idx])
		}
		idx++
	}
	if idx != 2 {
		t.Errorf("expected number of listed tuples to be 2, got %d", idx)
	}
}

func TestUnion(t *testing.T) {
	t.Run("schema mismatch", testUnionSchemaMismatch)
	t.Run("union", testUnion)
}

func testUnionSchemaMismatch(t *testing.T) {
	r1 := newRelation("foo")
	r2 := newRelation("bar")
	if _, err := relation.Union(r1, r2); err != relation.ErrSchemaMismatch {
		t.Errorf("expected error %v, got %v", relation.ErrSchemaMismatch, err)
	}
}

func testUnion(t *testing.T) {
	r1 := newRelation("foo")
	r2 := newRelation("foo")
	add(r1, map[string]any{"foo": 1})
	add(r2, map[string]any{"foo": 2})
	r, err := relation.Union(r1, r2)
	if err != nil {
		panic(err)
	}
	expect := []map[string]any{
		{"foo": 1},
		{"foo": 2},
	}
	idx := 0
	for t1 := range r.List() {
		if !maps.Equal(t1, expect[idx]) {
			t.Errorf("(%d) expected relation to contain %+v, got %+v", idx, expect[idx], t1)
		}
		idx++
	}
}

func newRelation(attributes ...string) *relation.Relation {
	r, err := relation.New(attributes...)
	if err != nil {
		panic(err)
	}
	return r
}

func add(r *relation.Relation, tuple map[string]any) *relation.Relation {
	if err := r.Add(tuple); err != nil {
		panic(err)
	}
	return r
}
