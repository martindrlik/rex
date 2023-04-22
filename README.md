# rex

Experimental relational NoSQL database. It is my playground for ideas and API will change over time. There is a lot more to do before it can be even considered interesting.

## Example

```
func Example() {
	rex.NewRelation().
		InsertOne(name("Jake"), bornYear(1980)).
		InsertOne(name("Lee"), bornYear(1979)).
		InsertOne(name("Kristen"), bornYear(1990)).
		Serialize(os.Stdout)
	// Output:
	// [{"bornYear": 1979, "name": "Lee"},
	// {"bornYear": 1980, "name": "Jake"},
	// {"bornYear": 1990, "name": "Kristen"}]
}
```