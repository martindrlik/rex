package table_test

import (
	"maps"
	"testing"

	"github.com/martindrlik/rex/table"
)

func TestNew(t *testing.T) {
	t.Run("missing schema no argument", func(t *testing.T) {
		if _, err := table.New(); err != table.ErrMissingSchema {
			t.Errorf("expected error %v, got %v", table.ErrMissingSchema, err)
		}
	})
	t.Run("missing schema empty slice", func(t *testing.T) {
		if _, err := table.New([]string{}...); err != table.ErrMissingSchema {
			t.Errorf("expected error %v, got %v", table.ErrMissingSchema, err)
		}
	})
}

func TestCount(t *testing.T) {
	t.Run("zero", func(t *testing.T) { expectCount(t, newTable(t, "foo"), 0) })
	t.Run("one", func(t *testing.T) { expectCount(t, add(t, newTable(t, "foo"), tup{"foo": 1}), 1) })
	t.Run("two", func(t *testing.T) {
		b1 := newTable(t, "foo")
		add(t, b1, tup{"foo": 1})
		add(t, b1, tup{"foo": 2})
		expectCount(t, b1, 2)
	})
	t.Run("three", func(t *testing.T) {
		b1 := newTable(t, "foo", "bar", "baz")
		add(t, b1, tup{"foo": 1})
		add(t, b1, tup{"foo": 2, "bar": 3})
		add(t, b1, tup{"foo": 4, "bar": 5, "baz": 6})
		expectCount(t, b1, 3)
	})
}

func TestAdd(t *testing.T) {
	t.Run("schema mismatch one", func(t *testing.T) {
		b1 := newTable(t, "foo")
		if err := b1.Add(tup{"bar": 1}); err != table.ErrSchemaMismatch {
			t.Errorf("expected error %v, got %v", table.ErrSchemaMismatch, err)
		}
	})
	t.Run("schema mismatch two", func(t *testing.T) {
		b1 := newTable(t, "foo")
		if err := b1.Add(tup{"foo": 1, "bar": 2}); err != table.ErrSchemaMismatch {
			t.Errorf("expected error %v, got %v", table.ErrSchemaMismatch, err)
		}
	})
	t.Run("missing schema", func(t *testing.T) {
		b1 := newTable(t, "foo")
		if err := b1.Add(tup{}); err != table.ErrMissingSchema {
			t.Errorf("expected error %v, got %v", table.ErrMissingSchema, err)
		}
	})
	t.Run("one", func(t *testing.T) {
		b1 := newTable(t, "foo")
		if err := b1.Add(tup{"foo": 1}); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
	t.Run("two", func(t *testing.T) {
		b1 := newTable(t, "foo", "bar")
		if err := b1.Add(tup{"foo": 1}); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if err := b1.Add(tup{"bar": 2}); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		expectCount(t, b1, 2)
	})
	t.Run("use copy", func(t *testing.T) {
		b1 := newTable(t, "foo", "bar")
		t1 := tup{"foo": 1, "bar": 2}
		add(t, b1, t1)
		t1["foo"] = 3
		expectTuples(t, b1, []tup{{"foo": 1, "bar": 2}})
	})
}

func TestContain(t *testing.T) {
	t.Run("not in empty table", func(t *testing.T) {
		b1 := newTable(t, "foo")
		if b1.Contain(tup{"foo": 1}) {
			t.Errorf("expected table to not contain tuple, got true")
		}
	})
	t.Run("different schema", func(t *testing.T) {
		b1 := newTable(t, "foo")
		add(t, b1, tup{"foo": 1})
		if b1.Contain(tup{"bar": 1}) {
			t.Errorf("expected table to not contain tuple, got true")
		}
	})
	t.Run("different value", func(t *testing.T) {
		b1 := newTable(t, "foo")
		add(t, b1, tup{"foo": 1})
		if b1.Contain(tup{"foo": 2}) {
			t.Errorf("expected table to not contain tuple, got true")
		}
	})
	t.Run("in table", func(t *testing.T) {
		b1 := newTable(t, "foo")
		add(t, b1, tup{"foo": 1})
		if !b1.Contain(tup{"foo": 1}) {
			t.Errorf("expected table to contain tuple, got false")
		}
	})
	t.Run("in table two", func(t *testing.T) {
		b1 := newTable(t, "foo", "bar")
		add(t, b1, tup{"foo": 1})
		add(t, b1, tup{"foo": 1, "bar": 2})
		if !b1.Contain(tup{"foo": 1, "bar": 2}) {
			t.Errorf("expected table to contain tuple, got false")
		}
	})
}

func TestList(t *testing.T) {
	t.Run("empty", func(t *testing.T) { expectTuples(t, newTable(t, "foo"), []tup{}) })
	t.Run("one", func(t *testing.T) {
		b1 := newTable(t, "foo")
		add(t, b1, tup{"foo": 1})
		expectTuples(t, b1, []tup{{"foo": 1}})
	})
	t.Run("two", func(t *testing.T) {
		b1 := newTable(t, "foo", "bar")
		add(t, b1, tup{"foo": 1})
		add(t, b1, tup{"foo": 1, "bar": 2})
		expectTuples(t, b1, []tup{{"foo": 1}, {"foo": 1, "bar": 2}})
	})
	t.Run("list one from two", func(t *testing.T) {
		b1 := newTable(t, "foo", "bar")
		add(t, b1, tup{"foo": 1})
		add(t, b1, tup{"foo": 1, "bar": 2})
		for t1 := range b1.List() {
			if !maps.Equal(t1, tup{"foo": 1}) {
				t.Errorf("expected %+v, got %+v", tup{"foo": 1}, t1)
			}
			break
		}
	})
}

func TestUnion(t *testing.T) {
	t.Run("schema mismatch one", func(t *testing.T) {
		b1 := newTable(t, "foo")
		b2 := newTable(t, "bar")
		if _, err := b1.Union(b2); err != table.ErrSchemaMismatch {
			t.Errorf("expected error %v, got %v", table.ErrSchemaMismatch, err)
		}
	})
	t.Run("schema mismatch two", func(t *testing.T) {
		b1 := newTable(t, "foo")
		b2 := newTable(t, "bar", "baz")
		if _, err := b1.Union(b2); err != table.ErrSchemaMismatch {
			t.Errorf("expected error %v, got %v", table.ErrSchemaMismatch, err)
		}
	})
	t.Run("one", func(t *testing.T) {
		b1 := newTable(t, "foo")
		b2 := newTable(t, "foo")
		add(t, b1, tup{"foo": 1})
		add(t, b2, tup{"foo": 2})
		expectTuples(t, union(t, b1, b2), []tup{{"foo": 1}, {"foo": 2}})
	})
}

func newTable(t *testing.T, aa ...string) *table.Table {
	t.Helper()
	b1, err := table.New(aa...)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	return b1
}

func add(t *testing.T, b1 *table.Table, t1 tup) *table.Table {
	t.Helper()
	if err := b1.Add(t1); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	return b1
}

func expectCount(t *testing.T, b1 *table.Table, n2 int) {
	t.Helper()
	if n1 := b1.Count(); n1 != n2 {
		t.Errorf("expected table to have %d tuple, got %d", n2, n1)
	}
}

func expectTuples(t *testing.T, b1 *table.Table, expect []tup) {
	t.Helper()
	if n1 := b1.Count(); n1 != len(expect) {
		t.Errorf("expected table to have %d tuples, got %d", len(expect), n1)
	}
	idx := 0
	for t1 := range b1.List() {
		if !maps.Equal(t1, expect[idx]) {
			t.Errorf("expected to have %+v, got %+v", expect[idx], t1)
		}
		idx++
	}
}

func union(t *testing.T, b1, b2 *table.Table) *table.Table {
	t.Helper()
	b3, err := b1.Union(b2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	return b3
}
