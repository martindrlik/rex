package table_test

import (
	"fmt"

	"github.com/martindrlik/rex/box"
	"github.com/martindrlik/rex/table"
)

func ExampleTable_Union() {
	movies2021 := table.New().Add(tuple(title("Dune: Part One"), year(2021), length(155)))
	movies2024 := table.New().Add(tuple(title("Dune: Part Two"), year(2024), length(166)))

	movies := movies2021.Union(movies2024)

	fmt.Println(box.Table([]string{"title", "year"}, movies.Tuples()...))

	// Output:
	// ┏━━━━━━━━━━━━━━━━┯━━━━━━┓
	// ┃ title          │ year ┃
	// ┠────────────────┼──────┨
	// ┃ Dune: Part One │ 2021 ┃
	// ┃ Dune: Part Two │ 2024 ┃
	// ┗━━━━━━━━━━━━━━━━┷━━━━━━┛
}
