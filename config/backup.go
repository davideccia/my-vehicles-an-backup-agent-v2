package config

import (
	"my_vehicles_an_backup_agent/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("backup", map[string]any{
		// Shared secret expected on the `x-api-key` header for every route
		// except /ping. Mirrors the NestJS reference server's API_KEY env var.
		"api_key": config.Env("API_KEY", ""),
		// The only accepted value for the `version` query/body field, sent by
		// the Android app as its Room DB_SCHEMA_VERSION. Any other value gets
		// a 409 Conflict, mirroring the NestJS reference server's SCHEMA_VERSION.
		"schema_version": "6",
	})
}
