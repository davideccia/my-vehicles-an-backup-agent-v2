package middleware

import (
	"crypto/subtle"

	"github.com/goravel/framework/contracts/http"

	"my_vehicles_an_backup_agent/app/facades"
)

// apiKey mirrors the NestJS reference server's ApiKeyGuard: every protected
// route requires an `x-api-key` header matching the configured secret,
// compared in constant time. /ping is intentionally excluded from this
// middleware's route group (see routes/api.go), matching @Public() there.
type apiKey struct{}

func ApiKey() http.Middleware {
	return &apiKey{}
}

func (m *apiKey) Signature() string {
	return "app:api_key"
}

func (m *apiKey) Handle(ctx http.Context) {
	expected := facades.Config().GetString("backup.api_key")
	received := ctx.Request().Header("x-api-key")

	if expected == "" || subtle.ConstantTimeCompare([]byte(received), []byte(expected)) != 1 {
		if err := ctx.Response().Status(http.StatusUnauthorized).Json(http.Json{
			"error": "unauthorized",
		}).Abort(); err != nil {
			panic(err)
		}
		return
	}

	ctx.Request().Next()
}
