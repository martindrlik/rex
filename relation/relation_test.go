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
	add(r, map[string]any{"foo": 1, "bar": 2})
	add(r, map[string]any{"foo": 3, "bar": 4})

	expect := []map[string]any{
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
	t.Run("schema mismatch", testUnionSchemaMismatch)
	t.Run("union", testUnion)
}

func testUnionSchemaMismatch(t *testing.T) {
	r1 := newRelation("foo")
	r2 := newRelation("bar")
	if _, err := r1.Union(r2); err != relation.ErrSchemaMismatch {
		t.Errorf("expected error %v, got %v", relation.ErrSchemaMismatch, err)
	}
}

func testUnion(t *testing.T) {
	r1 := newRelation("foo")
	r2 := newRelation("foo")
	add(r1, map[string]any{"foo": 1})
	add(r2, map[string]any{"foo": 2})
	r, err := r1.Union(r2)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expect := []map[string]any{
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
}

func TestDifference(t *testing.T) {
	t.Run("schema mismatch", testDifferenceSchemaMismatch)
	t.Run("difference", testDifference)
}

func testDifferenceSchemaMismatch(t *testing.T) {
	r1 := newRelation("foo")
	r2 := newRelation("bar")

	if _, err := r1.Difference(r2); err != relation.ErrSchemaMismatch {
		t.Errorf("expected error %v, got %v", relation.ErrSchemaMismatch, err)
	}
}

func testDifference(t *testing.T) {
	r1 := newRelation("foo")
	r2 := newRelation("foo")

	for i := 1; i <= 3; i++ {
		add(r1, map[string]any{"foo": i})
	}
	add(r2, map[string]any{"foo": 2})

	r3, err := r1.Difference(r2)
	if err != nil {
		panic(err)
	}

	expect := []map[string]any{
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
}

func TestCount(t *testing.T) {
	r := newRelation("foo")
	for i := 1; i <= 5; i++ {
		add(r, map[string]any{"foo": i})
	}
	if n := r.Count(); n != 5 {
		t.Errorf("expected relation to have 5 tuples, got %d", n)
	}
}

func TestNaturalJoin(t *testing.T) {
	t.Run("natural join", testNaturalJoin)
	t.Run("cartasian product", testNaturalJoinCartasianProduct)
}

func testNaturalJoin(t *testing.T) {
	r1 := newRelation("foo", "bar")
	r2 := newRelation("foo", "baz")

	add(r1, map[string]any{"foo": 1, "bar": 2})
	add(r1, map[string]any{"foo": 2, "bar": 3})
	add(r2, map[string]any{"foo": 1, "baz": 4})
	add(r2, map[string]any{"foo": 3, "baz": 5})

	r3 := r1.NaturalJoin(r2)

	expect := []map[string]any{
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
}

func testNaturalJoinCartasianProduct(t *testing.T) {
	r1 := newRelation("foo")
	r2 := newRelation("bar")

	add(r1, map[string]any{"foo": 1})
	add(r2, map[string]any{"bar": 1})
	add(r2, map[string]any{"bar": 2})

	r3 := r1.NaturalJoin(r2)

	expect := []map[string]any{
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
}

func TestRename(t *testing.T) {
	t.Run("rename", testRename)
	t.Run("duplicate attribute", testRenameDuplicateAttribute)
}

func testRenameDuplicateAttribute(t *testing.T) {
	r1 := newRelation("foo", "bar")
	if _, err := r1.Rename(map[string]string{"foo": "bar"}); err.Error() != relation.ErrDuplicateAttribute("bar").Error() {
		t.Errorf("expected error %v, got %v", relation.ErrDuplicateAttribute("bar").Error(), err.Error())
	}
}

func testRename(t *testing.T) {
	r1 := newRelation("bar", "baz")
	add(r1, map[string]any{"bar": 2, "baz": 1})
	type testCase struct {
		testName string
		input    map[string]string
		expect   []map[string]any
	}
	for tc := range func(yield func(testCase) bool) {
		yield(testCase{
			testName: "bar to foo",
			input:    map[string]string{"baz": "foo"},
			expect:   []map[string]any{{"foo": 1, "bar": 2}},
		})
	} {
		t.Run(tc.testName, func(t *testing.T) {
			r2, err := r1.Rename(tc.input)
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
