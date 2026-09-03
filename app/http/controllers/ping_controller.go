package controllers

import (
	"time"

	"github.com/goravel/framework/contracts/http"

	"my_vehicles_an_backup_agent/app/facades"
)

// PingController mirrors the NestJS reference server's public /ping route,
// used by the Android app both as a standalone connectivity check and as the
// first step of every backup/restore run.
type PingController struct{}

func NewPingController() *PingController {
	return &PingController{}
}

func (c *PingController) Ping(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"status":         "ok",
		"schema_version": facades.Config().GetString("backup.schema_version"),
		// ISO-8601 UTC with millisecond precision, matching the Kotlin client's
		// (opaque, but documented) expectation and the NestJS reference's
		// `new Date().toISOString()`.
		"timestamp": time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	})
}
