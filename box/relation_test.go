package box_test

import (
	"fmt"
	"slices"

	"github.com/martindrlik/rex/box"
)

func ExampleRelation() {
	fmt.Println(box.Relation(
		[]string{"title", "year"},
		slices.Values([]map[string]any{
			{"title": "Adventure Time", "year": 2010},
			{"title": "What We Do in the Shadows", "year": 2019},
			{"title": "The Last of Us"}})))

	fmt.Println(box.Relation([]string{"empty", "table"}, slices.Values([]map[string]any{})))

	// Output:
	// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━┓
	// ┃ title                     │ year ┃
	// ┠───────────────────────────┼──────┨
	// ┃ Adventure Time            │ 2010 ┃
	// ┃ What We Do in the Shadows │ 2019 ┃
	// ┃ The Last of Us            │ ?    ┃
	// ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━┛
	//
	// ┏━━━━━━━┯━━━━━━━┓
	// ┃ empty │ table ┃
	// ┗━━━━━━━┷━━━━━━━┛
}
