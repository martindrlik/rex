package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/martindrlik/rex/box"
	"github.com/martindrlik/rex/persist"
	"github.com/martindrlik/rex/table"
	"golang.org/x/exp/maps"
)

func main() {
	bind("union", "", func(a, b *table.Table) *table.Table { return a.Union(b) })
	bind("difference", "", func(a, b *table.Table) *table.Table { return a.Difference(b) })
	bind("intersection", "", func(a, b *table.Table) *table.Table { return a.Intersection(b) })
	bind("natural-join", "", func(a, b *table.Table) *table.Table { return a.NaturalJoin(b) })
	exec(parse(os.Args[1:]))
}

func exec(op string, tables []*table.Table, attributes []string, outputFormat string) {
	func(fn func([]*table.Table) []*table.Table) {
		for _, t := range fn(tables) {
			project(t, attributes, outputFormat)
		}
	}(binaryOp(op))
}

func binaryOp(op string) func([]*table.Table) []*table.Table {
	return aggr(func(a, b *table.Table) *table.Table {
		if desc, ok := ops[op]; ok {
			return desc.fn(a, b)
		}
		panic("unreachable")
	})
}

func aggr(fn func(a, b *table.Table) *table.Table) func([]*table.Table) []*table.Table {
	return func(tables []*table.Table) []*table.Table {
		result := tables[0]
		for _, t := range tables[1:] {
			result = fn(result, t)
		}
		return []*table.Table{result}
	}
}

func project(table *table.Table, attributes []string, outputFormat string) {
	if len(attributes) == 0 {
		attributes = maps.Keys(table.Schema())
		slices.Sort(attributes)
	}
	switch outputFormat {
	case "json":
		persist.StoreSchemaMode(os.Stdout, table.Project(attributes...))
	case "table":
		fmt.Println(box.Relation(attributes, table.Project(attributes...).Tuples()...))
	}
}

func parse(args []string) (string, []*table.Table, []string, string) {
	if len(args) < 2 {
		usage(errors.New("missing arguments"))
	}
	fs := flag.NewFlagSet("", flag.ExitOnError)
	var (
		schemalessFilenames = stringsFlag{}
		schemalessInlines   = stringsFlag{}

		schemaFilenames = stringsFlag{}
		schemaInlines   = stringsFlag{}

		outputFormat = fs.String("of", "table", "table or json")
	)
	fs.Var(&schemalessFilenames, "fa", "name of file that contains array of tuples")
	fs.Var(&schemalessInlines, "ia", "inline array of tuples")
	fs.Var(&schemaFilenames, "fs", "name of file that contains table object: schema and tuples")
	fs.Var(&schemaInlines, "is", "inline table object: schema and tuples")

	op := args[0]
	_, ok := ops[op]
	if !ok {
		usage(fmt.Errorf("unknown op: %s", op))
	}

	fs.Parse(args[1:])
	if len(schemalessFilenames) == 0 && len(schemaFilenames) == 0 && len(schemalessInlines) == 0 && len(schemaInlines) == 0 {
		usage(errors.New("missing table"))
	}

	tables := []*table.Table{}
	load := func(r io.Reader, fn func(io.Reader) (*table.Table, error)) error {
		t, err := fn(r)
		if err != nil {
			return err
		}
		tables = append(tables, t)
		return nil
	}
	loadFile := func(name string, fn func(io.Reader) (*table.Table, error)) error {
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()
		return load(f, fn)
	}
	loadFiles := func(filenames []string, fn func(io.Reader) (*table.Table, error)) {
		for _, name := range filenames {
			if err := loadFile(name, fn); err != nil {
				usage(fmt.Errorf("loading file: %w", err))
			}
		}
	}
	loadInline := func(inlines []string, fn func(io.Reader) (*table.Table, error)) {
		for _, inline := range inlines {
			t, err := fn(strings.NewReader(inline))
			if err != nil {
				usage(fmt.Errorf("loading inline %v: %w", inline, err))
			}
			tables = append(tables, t)
		}
	}

	loadFiles(schemalessFilenames, persist.Load)
	loadFiles(schemaFilenames, persist.LoadSchemaMode)

	loadInline(schemalessInlines, persist.Load)
	loadInline(schemaInlines, persist.LoadSchemaMode)

	return op, tables, fs.Args(), *outputFormat
}

type stringsFlag []string

func (s *stringsFlag) String() string {
	return fmt.Sprint(*s)
}

func (s *stringsFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func usage(err error) {
	if err != nil {
		fmt.Println("Error:")
		fmt.Printf("	%v\n", err)
	}
	fmt.Println("Usage:")
	fmt.Println("	rex <command> <input> <options> [attribute ...]")
	fmt.Println("Commands:")
	names := maps.Keys(ops)
	slices.Sort(names)
	for _, name := range names {
		fmt.Printf("	%s", name)
		desc := ops[name].desc
		if desc == "" {
			fmt.Println()
		} else {
			fmt.Printf("%s\n", desc)
		}
	}
	fmt.Println("Input:")
	fmt.Println("	-fa <file>   [-ta <file>   ...]: name of file that contains array of tuples")
	fmt.Println("	-ia <inline> [-ia <inline> ...]: inline array of tuples")
	fmt.Println("	-fs <file>   [-ts <file>   ...]: name of file that contains table object: schema and tuples")
	fmt.Println("	-is <inline> [-is <file>   ...]: inline table object: schema and tuples")
	fmt.Println("Options:")
	fmt.Println("	-of <format>: output format: table or json")

	fmt.Println("Note:")
	fmt.Println("	JSON is used as an input format")
	os.Exit(1)
}

type opDesc struct {
	desc string
	fn   func(a, b *table.Table) *table.Table
}

var ops = map[string]opDesc{}

func bind(name, desc string, fn func(a, b *table.Table) *table.Table) {
	ops[name] = opDesc{desc, fn}
}
