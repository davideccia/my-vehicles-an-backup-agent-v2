package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

// HealthController backs GET /healthcheck: an unauthenticated, unversioned
// liveness probe for the Docker HEALTHCHECK directive. Unlike /ping (which
// the Android app calls and which is version-gated to match the NestJS
// reference server), this endpoint answers to no contract but the
// container's own health check, so it skips both the API key and schema
// version middleware entirely.
type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) Check(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"status": "ok",
	})
}
