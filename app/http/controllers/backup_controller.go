package controllers

import (
	"fmt"
	"io"

	"github.com/goravel/framework/contracts/http"

	"my_vehicles_an_backup_agent/app/http/requests"
	"my_vehicles_an_backup_agent/app/models"
	"my_vehicles_an_backup_agent/app/services"
)

// BackupController implements the ten /backup/{table} routes (5 tables x
// POST+GET). Each handler is a thin wrapper around the generic postTable /
// getTable helpers below — Go methods can't be generic, so the table type is
// fixed per handler instead.
type BackupController struct{}

func NewBackupController() *BackupController {
	return &BackupController{}
}

func (c *BackupController) PostVehicles(ctx http.Context) http.Response {
	return postTable(ctx, requests.VehicleSpec, &models.Vehicle{})
}

func (c *BackupController) GetVehicles(ctx http.Context) http.Response {
	return getTable(ctx, &models.Vehicle{})
}

func (c *BackupController) PostVehicleRefuels(ctx http.Context) http.Response {
	return postTable(ctx, requests.VehicleRefuelSpec, &models.VehicleRefuel{})
}

func (c *BackupController) GetVehicleRefuels(ctx http.Context) http.Response {
	return getTable(ctx, &models.VehicleRefuel{})
}

func (c *BackupController) PostVehicleServiceTypes(ctx http.Context) http.Response {
	return postTable(ctx, requests.VehicleServiceTypeSpec, &models.VehicleServiceType{})
}

func (c *BackupController) GetVehicleServiceTypes(ctx http.Context) http.Response {
	return getTable(ctx, &models.VehicleServiceType{})
}

func (c *BackupController) PostVehicleServices(ctx http.Context) http.Response {
	return postTable(ctx, requests.VehicleServiceSpec, &models.VehicleService{})
}

func (c *BackupController) GetVehicleServices(ctx http.Context) http.Response {
	return getTable(ctx, &models.VehicleService{})
}

func (c *BackupController) PostVehicleServiceReminders(ctx http.Context) http.Response {
	return postTable(ctx, requests.VehicleServiceReminderSpec, &models.VehicleServiceReminder{})
}

func (c *BackupController) GetVehicleServiceReminders(ctx http.Context) http.Response {
	return getTable(ctx, &models.VehicleServiceReminder{})
}

// postTable validates the envelope's rows against spec, converts them to the
// concrete model type T, replaces the table's contents (services.SyncTable)
// and reports how many rows were kept vs. removed. Schema-version checking
// already happened in the SchemaVersion middleware.
func postTable[T any](ctx http.Context, spec requests.TableSpec, model *T) http.Response {
	body, err := io.ReadAll(ctx.Request().Origin().Body)
	if err != nil {
		return badRequest(ctx, "invalid_body")
	}

	env, err := requests.DecodeEnvelope(body)
	if err != nil {
		return badRequest(ctx, err.Error())
	}

	if err := requests.ValidateRows(env.Data, spec); err != nil {
		return badRequest(ctx, err.Error())
	}

	if dup := requests.DuplicateID(env.Data); dup != "" {
		return badRequest(ctx, fmt.Sprintf("duplicate id %q in payload", dup))
	}

	records, err := requests.Remarshal[T](env.Data)
	if err != nil {
		return badRequest(ctx, err.Error())
	}

	upserted, deleted, err := services.SyncTable(model, records, requests.IDs(env.Data))
	if err != nil {
		return internalError(ctx, err)
	}

	return ctx.Response().Status(http.StatusCreated).Json(http.Json{
		"upserted_count": upserted,
		"deleted_count":  deleted,
	})
}

// getTable returns the entire contents of a table. Schema-version checking
// already happened in the SchemaVersion middleware.
func getTable[T any](ctx http.Context, model *T) http.Response {
	records, err := services.FetchTable(model)
	if err != nil {
		return internalError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{
		"version": ctx.Request().Query("version"),
		"data":    records,
	})
}

func badRequest(ctx http.Context, message string) http.Response {
	return ctx.Response().Status(http.StatusBadRequest).Json(http.Json{
		"error": message,
	})
}
