package table_test

import (
	"testing"

	"github.com/martindrlik/rex/table"
)

func TestFilter(t *testing.T) {
	users, err := table.NewRelation(map[string]struct{}{
		"name": {},
		"age":  {},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	foo20 := map[string]any{"name": "Foo", "age": 20}
	bar21 := map[string]any{"name": "Bar", "age": 21}
	for _, tuple := range []map[string]any{foo20, bar21} {
		if err := users.Add(tuple); err != nil {
			t.Fatalf("expected no error, got %v, while adding tuple %+v", err, tuple)
		}
	}
	users21 := users.Filter(func(tuple map[string]any) bool {
		return tuple["age"] == 21
	})
	if users21.Has(foo20) {
		t.Errorf("expected users21 to have tuple %+v filtered out", foo20)
	}
	if !users21.Has(bar21) {
		t.Errorf("expected users21 to contain tuple %+v", bar21)
	}
}

func TestAreTuplesEqual(t *testing.T) {
	for name, testCase := range map[string]struct {
		a, b           map[string]any
		expectedResult bool
	}{
		"are equal": {
			a:              createUser("Foo", 20),
			b:              createUser("Foo", 20),
			expectedResult: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			if actual := table.AreTuplesEqual(testCase.a, testCase.b); actual != testCase.expectedResult {
				t.Errorf("testing tuple equality of following tuples %+v and %+v should return %v got %v", testCase.a, testCase.b, testCase.expectedResult, actual)
			}
		})
	}
}

func TestFirst(t *testing.T) {
	if !t.Run("make sure AreTuplesEqual is working", TestAreTuplesEqual) {
		t.Fatal("used AreTuplesEqual is not working as expected")
	}
}

func createUser(name string, age int) map[string]any {
	return map[string]any{"name": name, "age": age}
}
