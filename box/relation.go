package box

import (
	"fmt"
	"io"
	"iter"
	"strings"
	"unicode/utf8"
)

type relationBox struct {
	schema []string
	rows   []map[string]string
	max    map[string]int
}

func Relation(schema []string, tuples iter.Seq[map[string]any]) interface{ String() string } {
	rb := &relationBox{
		schema: schema,
		rows:   []map[string]string{},
		max:    map[string]int{},
	}
	for _, s := range rb.schema {
		rb.max[s] = utf8.RuneCountInString(s)
	}
	for t := range tuples {
		rb.addRow(t)
	}
	return rb
}

func (t *relationBox) addRow(tuple map[string]any) {
	str := func(v any) string { return fmt.Sprintf("%v", v) }
	row := map[string]string{}
	for k, v := range tuple {
		s := str(v)
		if l := utf8.RuneCountInString(s); t.max[k] < l {
			t.max[k] = l
		}
		row[k] = s
	}
	t.rows = append(t.rows, row)
}

func (t *relationBox) String() string {
	sb := &strings.Builder{}
	t.writeTop(sb)
	t.writeHeader(sb)
	if len(t.rows) > 0 {
		t.writeSeparator(sb)
		t.writeRows(sb)
	}
	t.writeBottom(sb)
	return sb.String()
}

func (t *relationBox) writeTop(w io.Writer) {
	// ┏━━━━━━┯━━━━━━┓
	t.writeRow(w, "┏", "┯", "┓", func(s string) string {
		return strings.Repeat("━", t.max[s]+2)
	})
}

func (t *relationBox) writeHeader(w io.Writer) {
	// ┃    x │    y ┃
	t.writeRow(w, "┃", "│", "┃", func(s string) string {
		return fmt.Sprintf(" %s ", t.pad(s, s))
	})
}

func (t *relationBox) writeSeparator(w io.Writer) {
	// ┠──────┼──────┨
	t.writeRow(w, "┠", "┼", "┨", func(s string) string {
		return strings.Repeat("─", t.max[s]+2)
	})
}

func (t *relationBox) writeRows(w io.Writer) {
	for _, row := range t.rows {
		// ┃ 2023 │ 2024 ┃
		t.writeRow(w, "┃", "│", "┃", func(s string) string {
			v, ok := row[s]
			return fmt.Sprintf(" %s ", t.pad(s, val(v, ok)))
		})
	}
}

func (t *relationBox) writeBottom(w io.Writer) {
	// ┗━━━━━━┷━━━━━━┛
	t.writeRow(w, "┗", "┷", "┛", func(s string) string {
		return strings.Repeat("━", t.max[s]+2)
	})
}

func val(v string, ok bool) string {
	if ok {
		return v
	}
	return "?"
}

func (t *relationBox) writeRow(w io.Writer, left, middle, right string, valueFunc func(string) string) {
	fmt.Fprint(w, left)
	for i, s := range t.schema {
		if i > 0 {
			fmt.Fprint(w, middle)
		}
		fmt.Fprint(w, valueFunc(s))
	}
	fmt.Fprintln(w, right)
}

func (bt *relationBox) pad(s, v string) string {
	return fmt.Sprintf("%s%s", v, strings.Repeat(" ", bt.max[s]-utf8.RuneCountInString(v)))
}
