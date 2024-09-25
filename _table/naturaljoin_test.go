package table_test

import (
	"fmt"

	"github.com/martindrlik/rex/box"
	"github.com/martindrlik/rex/table"
)

func ExampleTable_NaturalJoin() {
	matrix := func(m map[string]any) { title("The Matrix")(m) }
	dune := func(m map[string]any) { title("Dune: Part One")(m) }
	actor := func(name string) func(map[string]any) {
		return func(m map[string]any) {
			m["actor"] = name
		}
	}
	movies := table.New().Add(
		tuple(matrix, year(1999)),
		tuple(dune, year(2021)))
	actors := table.New().Add(
		tuple(actor("Keanu Reeves"), matrix),
		tuple(actor("Carrie-Anne Moss"), matrix),
		tuple(actor("Laurence Fishburne"), matrix),
		tuple(actor("Timothée Chalamet"), dune),
		tuple(actor("Rebecca Ferguson"), dune),
		tuple(actor("Zendaya"), dune))

	fmt.Println(box.Table([]string{"title", "year", "actor"}, movies.NaturalJoin(actors).Tuples()...))

	numbers := table.New().Add(
		T{"number": 1},
		T{"number": 2},
		T{"number": 3})
	letters := table.New().Add(
		T{"letter": "a"},
		T{"letter": "b"},
		T{"letter": "c"})

	fmt.Println("Cartesian product:")
	fmt.Println(box.Table([]string{"number", "letter"}, numbers.NaturalJoin(letters).Tuples()...))

	// Output:
	// ┏━━━━━━━━━━━━━━━━┯━━━━━━┯━━━━━━━━━━━━━━━━━━━━┓
	// ┃ title          │ year │ actor              ┃
	// ┠────────────────┼──────┼────────────────────┨
	// ┃ The Matrix     │ 1999 │ Keanu Reeves       ┃
	// ┃ The Matrix     │ 1999 │ Carrie-Anne Moss   ┃
	// ┃ The Matrix     │ 1999 │ Laurence Fishburne ┃
	// ┃ Dune: Part One │ 2021 │ Timothée Chalamet  ┃
	// ┃ Dune: Part One │ 2021 │ Rebecca Ferguson   ┃
	// ┃ Dune: Part One │ 2021 │ Zendaya            ┃
	// ┗━━━━━━━━━━━━━━━━┷━━━━━━┷━━━━━━━━━━━━━━━━━━━━┛
	//
	// Cartesian product:
	// ┏━━━━━━━━┯━━━━━━━━┓
	// ┃ number │ letter ┃
	// ┠────────┼────────┨
	// ┃ 1      │ a      ┃
	// ┃ 1      │ b      ┃
	// ┃ 1      │ c      ┃
	// ┃ 2      │ a      ┃
	// ┃ 2      │ b      ┃
	// ┃ 2      │ c      ┃
	// ┃ 3      │ a      ┃
	// ┃ 3      │ b      ┃
	// ┃ 3      │ c      ┃
	// ┗━━━━━━━━┷━━━━━━━━┛
}
