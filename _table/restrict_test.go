package table_test

import (
	"fmt"

	"github.com/martindrlik/rex/box"
	"github.com/martindrlik/rex/table"
)

func ExampleTable_Restrict() {
	movies := table.New().Add(
		movie("Die Hard", 1988),
		movie("The Matrix", 1999),
		movie("Guardians of the Galaxy", 2014),
		movie("Blade Runner: 2049", 2017),
		movie("Dune: Part One", 2021))

	year := func(f func(int) bool) func(tuple map[string]any) bool {
		return func(tuple map[string]any) bool {
			return f(tuple["year"].(int))
		}
	}

	fmt.Println(box.Table([]string{"title", "year"}, movies.Restrict(year(func(x int) bool { return x < 2000 })).Tuples()...))
	fmt.Println(box.Table([]string{"title", "year"}, movies.Restrict(year(func(x int) bool { return x > 2000 })).Tuples()...))

	// Output:
	// ┏━━━━━━━━━━━━┯━━━━━━┓
	// ┃ title      │ year ┃
	// ┠────────────┼──────┨
	// ┃ Die Hard   │ 1988 ┃
	// ┃ The Matrix │ 1999 ┃
	// ┗━━━━━━━━━━━━┷━━━━━━┛
	//
	// ┏━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━┓
	// ┃ title                   │ year ┃
	// ┠─────────────────────────┼──────┨
	// ┃ Guardians of the Galaxy │ 2014 ┃
	// ┃ Blade Runner: 2049      │ 2017 ┃
	// ┃ Dune: Part One          │ 2021 ┃
	// ┗━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━┛
}
