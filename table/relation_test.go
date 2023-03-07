package table_test

import (
	"maps"
	"testing"

	"github.com/martindrlik/rex/table"
)

func TestNewRelation(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		if _, err := table.NewRelation("foo"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
	t.Run("missing schema", func(t *testing.T) {
		if _, err := table.NewRelation(); err != table.ErrMissingSchema {
			t.Errorf("(no argument) expected error %v, got %v", table.ErrMissingSchema, err)
		}
		if _, err := table.NewRelation([]string{}...); err != table.ErrMissingSchema {
			t.Errorf("(empty slice) expected error %v, got %v", table.ErrMissingSchema, err)
		}
	})
}

func TestErrSchemaMismatch(t *testing.T) {
	r1 := newRelation(t, "foo")
	r2 := newRelation(t, "bar")

	type testCase struct {
		name string
		op   func(*table.Relation) (*table.Relation, error)
	}
	for tc := range func(yield func(testCase) bool) {
		yield(testCase{"union", r1.Union})
		yield(testCase{"difference", r1.Difference})
		yield(testCase{"intersection", r1.Intersection})
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.op(r2); err != table.ErrSchemaMismatch {
				t.Errorf("expected error %v, got %v", table.ErrSchemaMismatch, err)
			}
		})
	}
}

func TestRelationAdd(t *testing.T) {
	t.Run("count", TestCount)
	t.Run("one", func(t *testing.T) {
		r := newRelation(t, "foo", "bar", "baz")
		if err := r.Add(tup{"foo": 1, "bar": 2, "baz": 3}); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if n := r.Count(); n != 1 {
			t.Errorf("expected relation to have 1 tuple, got %d", n)
		}
	})
	t.Run("multiple", func(t *testing.T) {
		r := newRelation(t, "foo")
		if err := r.Add(tup{"foo": 1}); err != nil {
			t.Errorf("(adding foo: 1) expected no error, got %v", err)
		}
		if err := r.Add(tup{"foo": 2}); err != nil {
			t.Errorf("(adding foo: 2) expected no error, got %v", err)
		}
		if n := r.Count(); n != 2 {
			t.Errorf("expected relation to have 2 tuples, got %d", n)
		}
	})
	t.Run("already present", func(t *testing.T) {
		r := newRelation(t, "foo")
		if err := r.Add(tup{"foo": 1}); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if err := r.Add(tup{"foo": 1}); err != table.ErrAlreadyPresent {
			t.Errorf("expected error %v, got %v", table.ErrAlreadyPresent, err)
		}
	})
	t.Run("not matching schema", func(t *testing.T) {
		r := newRelation(t, "foo")
		t1 := tup{"bar": "baz"}
		t2 := tup{"foo": "1", "bar": "2"}
		t3 := tup{"foo": "bar"}
		if err := r.Add(t1); err != table.ErrSchemaMismatch {
			t.Errorf("(%+v) expected error %v, got %v", t1, table.ErrSchemaMismatch, err)
		}
		if err := r.Add(t2); err != table.ErrSchemaMismatch {
			t.Errorf("(%+v) expected error %v, got %v", t2, table.ErrSchemaMismatch, err)
		}
		if err := r.Add(t3); err != nil {
			t.Errorf("(%+v) expected to have matching schema, got error %v", t3, err)
		}
	})
}

func TestRelationContain(t *testing.T) {
	r := newRelation(t, "foo", "bar")
	if err := r.Add(tup{"foo": 1, "bar": 2}); err != nil {
		panic(err)
	}
	t.Run("true", func(t *testing.T) {
		if !r.Contain(tup{"foo": 1, "bar": 2}) {
			t.Errorf("expected relation to contain %+v", tup{"foo": 1, "bar": 2})
		}
	})
	t.Run("false different schema", func(t *testing.T) {
		if r.Contain(tup{"foo": 1}) {
			t.Error("expected relation to not contain tuple due to different schema")
		}
	})
	t.Run("false different value", func(t *testing.T) {
		if r.Contain(tup{"foo": 1, "bar": 3}) {
			t.Error("expected relation to not contain tuple due to different value")
		}
	})
	t.Run("false different values", func(t *testing.T) {
		if r.Contain(tup{"foo": 2, "bar": 3}) {
			t.Error("expected relation to not contain tuple due to different values")
		}
	})
}

func TestRelationList(t *testing.T) {
	r := newRelation(t, "foo", "bar")
	relationAdd(t, r, tup{"foo": 1, "bar": 2})
	relationAdd(t, r, tup{"foo": 3, "bar": 4})

	expect := []tup{
		{"foo": 1, "bar": 2},
		{"foo": 3, "bar": 4},
	}

	relationExpectCount(t, r, len(expect))

	idx := 0
	for t1 := range r.List() {
		if !maps.Equal(t1, expect[idx]) {
			t.Errorf("(%d) expected tuples %+v and %+v to be equal", idx, t1, expect[idx])
		}
		idx++
	}
}

func TestRelationUnion(t *testing.T) {
	t.Run("union", func(t *testing.T) {
		r1 := newRelation(t, "foo")
		r2 := newRelation(t, "foo")
		relationAdd(t, r1, tup{"foo": 1})
		relationAdd(t, r2, tup{"foo": 2})
		r, err := r1.Union(r2)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		expect := []tup{
			{"foo": 1},
			{"foo": 2},
		}

		relationExpectCount(t, r, len(expect))

		idx := 0
		for t1 := range r.List() {
			if !maps.Equal(t1, expect[idx]) {
				t.Errorf("(%d) expected relation to contain %+v, got %+v", idx, expect[idx], t1)
			}
			idx++
		}
	})
}

func TestRelationDifference(t *testing.T) {
	t.Run("difference", func(t *testing.T) {
		r1 := newRelation(t, "foo")
		r2 := newRelation(t, "foo")

		for i := 1; i <= 3; i++ {
			relationAdd(t, r1, tup{"foo": i})
		}
		relationAdd(t, r2, tup{"foo": 2})

		r3, err := r1.Difference(r2)
		if err != nil {
			panic(err)
		}

		expect := []tup{
			{"foo": 1},
			{"foo": 3}}

		relationExpectCount(t, r3, len(expect))

		idx := 0
		for t3 := range r3.List() {
			if !maps.Equal(t3, expect[idx]) {
				t.Errorf("expected to have %+v, got %+v", expect[idx], t3)
			}
			idx++
		}
	})
}

func TestRelationCount(t *testing.T) {
	r := newRelation(t, "foo")
	for i := 1; i <= 5; i++ {
		relationAdd(t, r, tup{"foo": i})
	}
	if n := r.Count(); n != 5 {
		t.Errorf("expected relation to have 5 tuples, got %d", n)
	}
}

func TestRelationNaturalJoin(t *testing.T) {
	t.Run("natural join", func(t *testing.T) {
		r1 := newRelation(t, "foo", "bar")
		r2 := newRelation(t, "foo", "baz")

		relationAdd(t, r1, tup{"foo": 1, "bar": 2})
		relationAdd(t, r1, tup{"foo": 2, "bar": 3})
		relationAdd(t, r2, tup{"foo": 1, "baz": 4})
		relationAdd(t, r2, tup{"foo": 3, "baz": 5})

		expectRelation(t, r1.NaturalJoin(r2), []tup{{"foo": 1, "bar": 2, "baz": 4}})
	})
	t.Run("cartasian product", func(t *testing.T) {
		r1 := newRelation(t, "foo")
		r2 := newRelation(t, "bar")

		relationAdd(t, r1, tup{"foo": 1})
		relationAdd(t, r2, tup{"bar": 1})
		relationAdd(t, r2, tup{"bar": 2})

		expectRelation(t, r1.NaturalJoin(r2), []tup{
			{"foo": 1, "bar": 1},
			{"foo": 1, "bar": 2},
		})
	})
}

func TestRelationRename(t *testing.T) {
	t.Run("rename", func(t *testing.T) {
		r1 := newRelation(t, "bar", "baz")
		relationAdd(t, r1, tup{"bar": 2, "baz": 1})
		type testCase struct {
			testName string
			renameFn func(string) string
			expect   []map[string]any
		}
		for tc := range func(yield func(testCase) bool) {
			yield(testCase{
				testName: "bar to foo",
				renameFn: func(old string) string {
					if old == "baz" {
						return "foo"
					}
					return old
				},
				expect: []tup{{"foo": 1, "bar": 2}},
			})
		} {
			t.Run(tc.testName, func(t *testing.T) {
				r2, err := r1.Rename(tc.renameFn)
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				relationExpectCount(t, r2, len(tc.expect))

				idx := 0
				for t2 := range r2.List() {
					if !maps.Equal(t2, tc.expect[idx]) {
						t.Errorf("expected to have %+v, got %+v", tc.expect[idx], t2)
					}
					idx++
				}
			})
		}
	})
	t.Run("duplicate attribute", func(t *testing.T) {
		r1 := newRelation(t, "foo", "bar")
		if _, err := r1.Rename(func(old string) (new string) {
			if old == "foo" {
				return "bar"
			}
			return old
		}); err.Error() != table.ErrDuplicateAttribute("bar").Error() {
			t.Errorf("expected error %v, got %v", table.ErrDuplicateAttribute("bar").Error(), err.Error())
		}
	})
}

func TestRelationIntersection(t *testing.T) {
	t.Run("intersection", func(t *testing.T) {
		r1 := newRelation(t, "foo")
		r2 := newRelation(t, "foo")

		relationAdd(t, r1, tup{"foo": 1})
		relationAdd(t, r1, tup{"foo": 2})
		relationAdd(t, r2, tup{"foo": 2})
		relationAdd(t, r2, tup{"foo": 3})

		r3, err := r1.Intersection(r2)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		expectRelation(t, r3, []tup{{"foo": 2}})
	})
}

func TestRelationProject(t *testing.T) {
	t.Run("missing schema", func(t *testing.T) {
		r1 := newRelation(t, "foo", "bar", "baz")
		relationAdd(t, r1, tup{"foo": 1, "bar": 2, "baz": 3})
		if _, err := r1.Project(); err != table.ErrMissingSchema {
			t.Errorf("(no argument) expected error %v, got %v", table.ErrMissingSchema, err)
		}
		if _, err := r1.Project([]string{}...); err != table.ErrMissingSchema {
			t.Errorf("(empty slice) expected error %v, got %v", table.ErrMissingSchema, err)
		}
	})
	t.Run("schema mismatch", func(t *testing.T) {
		r1 := newRelation(t, "foo", "bar", "baz")
		relationAdd(t, r1, tup{"foo": 1, "bar": 2, "baz": 3})
		if _, err := r1.Project("pub"); err != table.ErrSchemaMismatch {
			t.Errorf("expected error %v, got %v", table.ErrSchemaMismatch, err)
		}
	})
	t.Run("project", func(t *testing.T) {
		r1 := newRelation(t, "foo", "bar", "baz")
		relationAdd(t, r1, tup{"foo": 1, "bar": 2, "baz": 3})
		r2, err := r1.Project("bar")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		expectRelation(t, r2, []tup{{"bar": 2}})
	})
}

func newRelation(t *testing.T, aa ...string) *table.Relation {
	t.Helper()
	r, err := table.NewRelation(aa...)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	return r
}

func relationAdd(t *testing.T, r *table.Relation, tuple map[string]any) *table.Relation {
	t.Helper()
	if err := r.Add(tuple); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	return r
}

func relationExpectCount(t *testing.T, r *table.Relation, expectedCount int) {
	t.Helper()
	if n := r.Count(); n != expectedCount {
		t.Errorf("expected relation to have %d tuples, got %d", expectedCount, n)
	}
}

func expectRelation(t *testing.T, r *table.Relation, expect []tup) {
	t.Helper()
	if n := r.Count(); n != len(expect) {
		t.Errorf("expected relation to have %d tuples, got %d", len(expect), n)
	}

	idx := 0
	for t1 := range r.List() {
		if !maps.Equal(t1, expect[idx]) {
			t.Errorf("expected to have %+v, got %+v", expect[idx], t1)
		}
		idx++
	}
}
