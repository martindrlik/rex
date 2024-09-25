package table_test

import (
	"fmt"

	"github.com/martindrlik/rex/box"
	"github.com/martindrlik/rex/table"
)

func ExampleTable_Intersection() {
	u := table.New().Add(
		tuple(title("Dune: Part One")),
		movie("Dune: Part Two", 2024))
	v := table.New().Add(tuple(title("Dune: Part One")))

	fmt.Println(box.Table([]string{"title"}, u.Intersection(v).Tuples()...))

	// Output:
	// ┏━━━━━━━━━━━━━━━━┓
	// ┃ title          ┃
	// ┠────────────────┨
	// ┃ Dune: Part One ┃
	// ┗━━━━━━━━━━━━━━━━┛
}
