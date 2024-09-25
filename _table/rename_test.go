package table_test

import (
	"fmt"

	"github.com/martindrlik/rex/box"
	"github.com/martindrlik/rex/table"
)

func ExampleTable_Rename() {
	movies := table.New().Add(
		movie("The Matrix", 1999),
		movie("Dune: Part One", 2021))
	movies = movies.Rename(map[string]string{"title": "movie_title", "year": "released"})
	fmt.Println(box.Table(
		[]string{"movie_title", "released"},
		movies.Tuples()...))

	// Output:
	// ┏━━━━━━━━━━━━━━━━┯━━━━━━━━━━┓
	// ┃ movie_title    │ released ┃
	// ┠────────────────┼──────────┨
	// ┃ The Matrix     │ 1999     ┃
	// ┃ Dune: Part One │ 2021     ┃
	// ┗━━━━━━━━━━━━━━━━┷━━━━━━━━━━┛
}
