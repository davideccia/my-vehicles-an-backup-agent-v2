// Package requests holds the request-body parsing and validation shared by
// the /backup/{table} endpoints. Every table is posted as the same envelope:
// `{"version": "...", "data": [...]}`, so the envelope decoding and per-row
// checks (unknown fields, required fields, enum values, duplicate ids) live
// here once instead of being repeated for each of the five tables.
package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Envelope is the generic POST /backup/{table} request body.
type Envelope struct {
	Version string           `json:"version"`
	Data    []map[string]any `json:"data"`
}

// DecodeEnvelope parses body strictly: an unknown top-level key (anything
// besides "version" and "data") is rejected, mirroring the NestJS reference
// server's `forbidNonWhitelisted` validation pipe.
func DecodeEnvelope(body []byte) (Envelope, error) {
	var env Envelope
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&env); err != nil {
		return Envelope{}, fmt.Errorf("invalid request body: %w", err)
	}
	return env, nil
}

// Remarshal converts validated row maps into the concrete model type T,
// rejecting any field whose JSON type doesn't match T's (e.g. a string where
// an int is expected). Row-level checks (required/enum/unknown-key) must run
// before this, since a type mismatch here only reports the first offending
// row via the standard json error, not a full per-row report.
func Remarshal[T any](rows []map[string]any) ([]T, error) {
	buf, err := json.Marshal(rows)
	if err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	dec := json.NewDecoder(bytes.NewReader(buf))
	dec.DisallowUnknownFields()

	var out []T
	if err := dec.Decode(&out); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return out, nil
}
