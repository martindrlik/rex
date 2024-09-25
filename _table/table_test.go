package table_test

import (
	"fmt"
	"testing"

	"github.com/martindrlik/rex/box"
	"github.com/martindrlik/rex/table"
)

type T = map[string]any

func ExampleTable() {
	movies := table.New("title", "year", "length").Add(
		movie("The Matrix", 1999),
		tuple(title("Blade Runner: 2049"), year(2017), length(164)),
		tuple(title("Dune: Part One"), year(2021), length(155)))

	fmt.Println(box.Table(movies.SchemaOrder(), movies.Tuples()...))

	// Output:
	// ┏━━━━━━━━━━━━━━━━━━━━┯━━━━━━┯━━━━━━━━┓
	// ┃ title              │ year │ length ┃
	// ┠────────────────────┼──────┼────────┨
	// ┃ The Matrix         │ 1999 │ ?      ┃
	// ┃ Blade Runner: 2049 │ 2017 │ 164    ┃
	// ┃ Dune: Part One     │ 2021 │ 155    ┃
	// ┗━━━━━━━━━━━━━━━━━━━━┷━━━━━━┷━━━━━━━━┛
}

func ExampleTable_Add() {
	movies := table.New().Add(
		movie("The Matrix", 1999),
		movie("The Matrix", 1999)) // duplicate

	fmt.Println(box.Table([]string{"title", "year"}, movies.Tuples()...))

	// Output:
	// ┏━━━━━━━━━━━━┯━━━━━━┓
	// ┃ title      │ year ┃
	// ┠────────────┼──────┨
	// ┃ The Matrix │ 1999 ┃
	// ┗━━━━━━━━━━━━┷━━━━━━┛
}

func ExampleTable_Remove() {
	theMatrix := movie("The Matrix", 1999)
	movies := table.New().Add(theMatrix, movie("Dune: Part Two", 2023))
	movies = movies.Remove(theMatrix)
	fmt.Println(box.Table([]string{"title", "year"}, movies.Tuples()...))

	// Output:
	// ┏━━━━━━━━━━━━━━━━┯━━━━━━┓
	// ┃ title          │ year ┃
	// ┠────────────────┼──────┨
	// ┃ Dune: Part Two │ 2023 ┃
	// ┗━━━━━━━━━━━━━━━━┷━━━━━━┛
}

func TestContains(t *testing.T) {
	matrixMovie := movie("The Matrix", 1999)
	movies := table.New().Add(matrixMovie)
	moviesBox := box.Table([]string{"title", "year"}, movies.Tuples()...)
	if !movies.Tuples().Contains(matrixMovie) {
		t.Errorf(
			"\nexpected\n%v\nto contain\n%v",
			moviesBox,
			matrixMovie)
	}

	matrixWithLength := tuple(title("The Matrix"), year(1999), length(136))
	if movies.Tuples().Contains(matrixWithLength) {
		t.Errorf(
			"\nexpected\n%v\nnot to contain\n%v",
			moviesBox,
			matrixWithLength)
	}
}

func movie(t string, y int) map[string]any {
	return tuple(title(t), year(y))
}

func title(title string) func(map[string]any) {
	return func(tuple map[string]any) {
		tuple["title"] = title
	}
}

func year(year int) func(map[string]any) {
	return func(tuple map[string]any) {
		tuple["year"] = year
	}
}

func length(length int) func(map[string]any) {
	return func(tuple map[string]any) {
		tuple["length"] = length
	}
}

func tuple(set ...func(map[string]any)) map[string]any {
	x := map[string]any{}
	for _, f := range set {
		f(x)
	}
	return x
}
