package rex_test

import (
	"testing"

	"github.com/martindrlik/rex"
)

func TestIntersect(t *testing.T) {
	r := rex.NewRelation().
		InsertOne(name("Harry")).
		InsertOne(name("Sally")).
		InsertOne(name("George"))
	s := rex.NewRelation().
		InsertOne(name("George"))
	expected := rex.NewRelation().
		InsertOne(name("George"))
	if !expected.Equals(r.Intersect(s)) {
		t.Error("expected equal after intersect")
	}
}