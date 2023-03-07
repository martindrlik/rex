package table_test

import (
	"fmt"

	"github.com/martindrlik/rex/box"
	"github.com/martindrlik/rex/table"
)

func ExampleTable_Project() {
	movies := table.New().Add(
		movie("The Matrix", 1999),
		movie("Dune: Part One", 2021),
		tuple(title("Blade Runner: 2049"), year(2017), length(164)))
	titles := movies.Project("title")
	fmt.Println(box.Table(titles.SchemaOrder(), titles.Tuples()...))

	// Output:
	// ┏━━━━━━━━━━━━━━━━━━━━┓
	// ┃ title              ┃
	// ┠────────────────────┨
	// ┃ The Matrix         ┃
	// ┃ Dune: Part One     ┃
	// ┃ Blade Runner: 2049 ┃
	// ┗━━━━━━━━━━━━━━━━━━━━┛
}
