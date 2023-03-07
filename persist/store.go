package persist

import (
	"encoding/json"
	"io"

	"github.com/martindrlik/rex/table"
)

func StoreSchemaMode(w io.Writer, t *table.Table) error {
	enc := json.NewEncoder(w)
	return enc.Encode(struct {
		Schema []string         `json:"schema"`
		Tuples []map[string]any `json:"tuples"`
	}{
		Schema: t.SchemaOrder(),
		Tuples: t.Tuples(),
	})
}
