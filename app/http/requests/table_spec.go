package requests

import (
	"fmt"
	"slices"
)

// FieldSpec describes one allowed field of a backup row: whether it must be
// present and non-null, and — for the two enum-like columns the NestJS DTOs
// validated with @IsEnum — the values it may take.
type FieldSpec struct {
	Name     string
	Required bool
	Enum     []string
}

// TableSpec is the full set of fields a table's backup rows may contain.
type TableSpec struct {
	Fields []FieldSpec
}

func (s TableSpec) allowedKeys() map[string]struct{} {
	m := make(map[string]struct{}, len(s.Fields))
	for _, f := range s.Fields {
		m[f.Name] = struct{}{}
	}
	return m
}

// ValidateRows checks every row against spec: unknown keys, required/non-null
// fields, and enum membership. It returns the first violation found, or nil
// if every row passes. JSON type-checking (int vs string, etc.) is left to
// Remarshal, which runs afterward.
func ValidateRows(rows []map[string]any, spec TableSpec) error {
	allowed := spec.allowedKeys()

	for i, row := range rows {
		for key := range row {
			if _, ok := allowed[key]; !ok {
				return fmt.Errorf("data[%d]: unknown field %q", i, key)
			}
		}

		for _, field := range spec.Fields {
			val, present := row[field.Name]

			if field.Required && (!present || val == nil) {
				return fmt.Errorf("data[%d]: %q is required", i, field.Name)
			}

			if len(field.Enum) == 0 || !present || val == nil {
				continue
			}

			str, ok := val.(string)
			if !ok || !slices.Contains(field.Enum, str) {
				return fmt.Errorf("data[%d]: %q must be one of %v", i, field.Name, field.Enum)
			}
		}
	}

	return nil
}

// DuplicateID returns the first id that appears more than once in rows, or ""
// if every id is unique. A duplicate would violate the primary key once
// inserted, surfacing as an opaque 500 instead of a clear 400.
func DuplicateID(rows []map[string]any) string {
	seen := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		id, _ := row["id"].(string)
		if _, ok := seen[id]; ok {
			return id
		}
		seen[id] = struct{}{}
	}
	return ""
}

// IDs extracts the "id" field of every row, in order. Rows without a string
// id yield "" for that position — ValidateRows rejects that before this is
// ever called in practice, since "id" is Required on every table.
func IDs(rows []map[string]any) []string {
	ids := make([]string, len(rows))
	for i, row := range rows {
		ids[i], _ = row["id"].(string)
	}
	return ids
}
