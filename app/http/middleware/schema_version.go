package middleware

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/goravel/framework/contracts/http"

	"my_vehicles_an_backup_agent/app/facades"
)

// schemaVersion mirrors the NestJS reference server's assertSchemaVersion:
// every request (GET and POST alike) must carry the exact expected `version`
// value, or it is rejected with 409 Conflict before touching the database.
//
// On GET requests version travels as a query parameter. On POST requests it
// is a field in the JSON body, so this middleware peeks at the body just for
// that field, then restores it (the request body is a one-shot reader) so
// the controller can still decode it in full.
type schemaVersion struct{}

func SchemaVersion() http.Middleware {
	return &schemaVersion{}
}

func (m *schemaVersion) Signature() string {
	return "app:schema_version"
}

func (m *schemaVersion) Handle(ctx http.Context) {
	expected := facades.Config().GetString("backup.schema_version")
	received := ctx.Request().Query("version")

	if ctx.Request().Method() == http.MethodPost {
		origin := ctx.Request().Origin()
		body, err := io.ReadAll(origin.Body)
		if err != nil {
			if abortErr := ctx.Response().Status(http.StatusBadRequest).Json(http.Json{
				"error": "invalid_body",
			}).Abort(); abortErr != nil {
				panic(abortErr)
			}
			return
		}
		origin.Body = io.NopCloser(bytes.NewReader(body))

		var peek struct {
			Version string `json:"version"`
		}
		// A malformed body is reported by the controller's strict decode;
		// here a failed peek just leaves received empty, which fails the
		// version check below.
		_ = json.Unmarshal(body, &peek)
		received = peek.Version
	}

	if received != expected {
		if err := ctx.Response().Status(http.StatusConflict).Json(http.Json{
			"error":    "version_mismatch",
			"expected": expected,
			"received": received,
		}).Abort(); err != nil {
			panic(err)
		}
		return
	}

	ctx.Request().Next()
}
