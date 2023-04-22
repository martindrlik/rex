# rex

Experimental relational NoSQL database. It is my playground for ideas and API will change over time. There is a lot more to do before it can be even considered interesting.

## Example

```
func Example() {
	names := rex.NewRelation().InsertManyJson(strings.NewReader(`[
		{"id": 1, "name": "Lee"},
		{"id": 2, "name": "Jake"},
		{"id": 3, "name": "Kristen"}
	]`))
	years := rex.NewRelation().InsertManyJson(strings.NewReader(`[
		{"id": 1, "bornYear": 1979},
		{"id": 2, "bornYear": 1980},
		{"id": 3, "bornYear": 1990}
	]`))
	names.NaturalJoin(years).Project("bornYear", "name").Serialize(os.Stdout)
	// Output:
	// [{"bornYear": 1979, "name": "Lee"},
	// {"bornYear": 1980, "name": "Jake"},
	// {"bornYear": 1990, "name": "Kristen"}]
}
```