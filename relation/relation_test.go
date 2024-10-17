package relation_test

import (
	"maps"
	"testing"

	"github.com/martindrlik/rex/relation"
)

func TestNew(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		if _, err := relation.New("foo"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
	t.Run("missing schema", func(t *testing.T) {
		if _, err := relation.New(); err != relation.ErrMissingSchema {
			t.Errorf("(no argument) expected error %v, got %v", relation.ErrMissingSchema, err)
		}
		if _, err := relation.New([]string{}...); err != relation.ErrMissingSchema {
			t.Errorf("(empty slice) expected error %v, got %v", relation.ErrMissingSchema, err)
		}
	})
}

func TestErrSchemaMismatch(t *testing.T) {
	r1 := newRelation("foo")
	r2 := newRelation("bar")

	type testCase struct {
		name string
		op   func(*relation.Relation) (*relation.Relation, error)
	}
	for tc := range func(yield func(testCase) bool) {
		yield(testCase{"union", r1.Union})
		yield(testCase{"difference", r1.Difference})
		yield(testCase{"intersection", r1.Intersection})
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.op(r2); err != relation.ErrSchemaMismatch {
				t.Errorf("expected error %v, got %v", relation.ErrSchemaMismatch, err)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	t.Run("count", TestCount)
	t.Run("one", func(t *testing.T) {
		r := newRelation("foo", "bar", "baz")
		if err := r.Add(tup{"foo": 1, "bar": 2, "baz": 3}); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if n := r.Count(); n != 1 {
			t.Errorf("expected relation to have 1 tuple, got %d", n)
		}
	})
	t.Run("multiple", func(t *testing.T) {
		r := newRelation("foo")
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
		r := newRelation("foo")
		if err := r.Add(tup{"foo": 1}); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if err := r.Add(tup{"foo": 1}); err != relation.ErrAlreadyPresent {
			t.Errorf("expected error %v, got %v", relation.ErrAlreadyPresent, err)
		}
	})
	t.Run("not matching schema", func(t *testing.T) {
		r := newRelation("foo")
		t1 := tup{"bar": "baz"}
		t2 := tup{"foo": "1", "bar": "2"}
		t3 := tup{"foo": "bar"}
		if err := r.Add(t1); err != relation.ErrSchemaMismatch {
			t.Errorf("(%+v) expected error %v, got %v", t1, relation.ErrSchemaMismatch, err)
		}
		if err := r.Add(t2); err != relation.ErrSchemaMismatch {
			t.Errorf("(%+v) expected error %v, got %v", t2, relation.ErrSchemaMismatch, err)
		}
		if err := r.Add(t3); err != nil {
			t.Errorf("(%+v) expected to have matching schema, got error %v", t3, err)
		}
	})
}

func TestContain(t *testing.T) {
	r := newRelation("foo", "bar")
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

func TestList(t *testing.T) {
	r := newRelation("foo", "bar")
	add(r, tup{"foo": 1, "bar": 2})
	add(r, tup{"foo": 3, "bar": 4})

	expect := []tup{
		{"foo": 1, "bar": 2},
		{"foo": 3, "bar": 4},
	}

	expectCount(r, len(expect), t)

	idx := 0
	for t1 := range r.List() {
		if !maps.Equal(t1, expect[idx]) {
			t.Errorf("(%d) expected tuples %+v and %+v to be equal", idx, t1, expect[idx])
		}
		idx++
	}
}

func TestUnion(t *testing.T) {
	t.Run("union", func(t *testing.T) {
		r1 := newRelation("foo")
		r2 := newRelation("foo")
		add(r1, tup{"foo": 1})
		add(r2, tup{"foo": 2})
		r, err := r1.Union(r2)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		expect := []tup{
			{"foo": 1},
			{"foo": 2},
		}

		expectCount(r, len(expect), t)

		idx := 0
		for t1 := range r.List() {
			if !maps.Equal(t1, expect[idx]) {
				t.Errorf("(%d) expected relation to contain %+v, got %+v", idx, expect[idx], t1)
			}
			idx++
		}
	})
}

func TestDifference(t *testing.T) {
	t.Run("difference", func(t *testing.T) {
		r1 := newRelation("foo")
		r2 := newRelation("foo")

		for i := 1; i <= 3; i++ {
			add(r1, tup{"foo": i})
		}
		add(r2, tup{"foo": 2})

		r3, err := r1.Difference(r2)
		if err != nil {
			panic(err)
		}

		expect := []tup{
			{"foo": 1},
			{"foo": 3}}

		expectCount(r3, len(expect), t)

		idx := 0
		for t3 := range r3.List() {
			if !maps.Equal(t3, expect[idx]) {
				t.Errorf("expected to have %+v, got %+v", expect[idx], t3)
			}
			idx++
		}
	})
}

func TestCount(t *testing.T) {
	r := newRelation("foo")
	for i := 1; i <= 5; i++ {
		add(r, tup{"foo": i})
	}
	if n := r.Count(); n != 5 {
		t.Errorf("expected relation to have 5 tuples, got %d", n)
	}
}

func TestNaturalJoin(t *testing.T) {
	t.Run("natural join", func(t *testing.T) {
		r1 := newRelation("foo", "bar")
		r2 := newRelation("foo", "baz")

		add(r1, tup{"foo": 1, "bar": 2})
		add(r1, tup{"foo": 2, "bar": 3})
		add(r2, tup{"foo": 1, "baz": 4})
		add(r2, tup{"foo": 3, "baz": 5})

		r3 := r1.NaturalJoin(r2)

		expect := []tup{
			{"foo": 1, "bar": 2, "baz": 4},
		}

		expectCount(r3, len(expect), t)

		idx := 0
		for t3 := range r3.List() {
			if !maps.Equal(t3, expect[idx]) {
				t.Errorf("expected to have %+v, got %+v", expect[idx], t3)
			}
			idx++
		}
	})
	t.Run("cartasian product", func(t *testing.T) {
		r1 := newRelation("foo")
		r2 := newRelation("bar")

		add(r1, tup{"foo": 1})
		add(r2, tup{"bar": 1})
		add(r2, tup{"bar": 2})

		r3 := r1.NaturalJoin(r2)

		expect := []tup{
			{"foo": 1, "bar": 1},
			{"foo": 1, "bar": 2},
		}

		expectCount(r3, len(expect), t)

		idx := 0
		for t3 := range r3.List() {
			if !maps.Equal(t3, expect[idx]) {
				t.Errorf("expected to have %+v, got %+v", expect[idx], t3)
			}
			idx++
		}
	})
}

func TestRename(t *testing.T) {
	t.Run("rename", func(t *testing.T) {
		r1 := newRelation("bar", "baz")
		add(r1, tup{"bar": 2, "baz": 1})
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

				expectCount(r2, len(tc.expect), t)

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
		r1 := newRelation("foo", "bar")
		if _, err := r1.Rename(func(old string) (new string) {
			if old == "foo" {
				return "bar"
			}
			return old
		}); err.Error() != relation.ErrDuplicateAttribute("bar").Error() {
			t.Errorf("expected error %v, got %v", relation.ErrDuplicateAttribute("bar").Error(), err.Error())
		}
	})
}

func TestIntersection(t *testing.T) {
	t.Run("intersection", func(t *testing.T) {
		r1 := newRelation("foo")
		r2 := newRelation("foo")

		add(r1, tup{"foo": 1})
		add(r1, tup{"foo": 2})
		add(r2, tup{"foo": 2})
		add(r2, tup{"foo": 3})

		r3, err := r1.Intersection(r2)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		expect := []tup{{"foo": 2}}

		expectCount(r3, len(expect), t)

		idx := 0
		for t3 := range r3.List() {
			if !maps.Equal(t3, expect[idx]) {
				t.Errorf("expected to have %+v, got %+v", expect[idx], t3)
			}
			idx++
		}
	})
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

func expectCount(r *relation.Relation, expectedCount int, t *testing.T) {
	if n := r.Count(); n != expectedCount {
		t.Errorf("expected relation to have %d tuples, got %d", expectedCount, n)
	}
}

type tup = map[string]any
