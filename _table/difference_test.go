package table_test

import (
	"fmt"

	"github.com/martindrlik/rex/box"
	"github.com/martindrlik/rex/table"
)

func ExampleTable_Difference() {
	available := table.New().Add(
		tuple(title("Dune: Part One")),
		tuple(title("Dune: Part Two"), year(2024)))

	fmt.Println(box.Table([]string{"title"}, available.Difference(table.New().Add(tuple(title("Dune: Part One")))).Tuples()...))
	fmt.Println(box.Table([]string{"title"}, available.Difference(table.New().Add(tuple(title("Dune: Part Two"), year(2024)))).Tuples()...))

	// Output:
	// ┏━━━━━━━━━━━━━━━━┓
	// ┃ title          ┃
	// ┠────────────────┨
	// ┃ Dune: Part One ┃
	// ┃ Dune: Part Two ┃
	// ┗━━━━━━━━━━━━━━━━┛
	//
	// ┏━━━━━━━━━━━━━━━━┓
	// ┃ title          ┃
	// ┠────────────────┨
	// ┃ Dune: Part One ┃
	// ┗━━━━━━━━━━━━━━━━┛
}
