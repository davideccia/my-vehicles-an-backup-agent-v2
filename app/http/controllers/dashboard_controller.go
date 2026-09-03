package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"my_vehicles_an_backup_agent/app/facades"
	"my_vehicles_an_backup_agent/app/models"
	"my_vehicles_an_backup_agent/app/services"
)

// DashboardController mirrors the NestJS reference server's
// GET /dashboard/stats: five independent row counts, wrapped in a `data`
// envelope.
type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

func (c *DashboardController) Stats(ctx http.Context) http.Response {
	vehicles, err := services.CountTable(&models.Vehicle{})
	if err != nil {
		return internalError(ctx, err)
	}

	refuels, err := services.CountTable(&models.VehicleRefuel{})
	if err != nil {
		return internalError(ctx, err)
	}

	serviceCount, err := services.CountTable(&models.VehicleService{})
	if err != nil {
		return internalError(ctx, err)
	}

	serviceTypes, err := services.CountTable(&models.VehicleServiceType{})
	if err != nil {
		return internalError(ctx, err)
	}

	serviceReminders, err := services.CountTable(&models.VehicleServiceReminder{})
	if err != nil {
		return internalError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{
		"data": http.Json{
			"vehicles":          vehicles,
			"refuels":           refuels,
			"services":          serviceCount,
			"service_types":     serviceTypes,
			"service_reminders": serviceReminders,
		},
	})
}

func internalError(ctx http.Context, err error) http.Response {
	facades.Log().Error(err)
	return ctx.Response().Status(http.StatusInternalServerError).Json(http.Json{
		"error": "internal_error",
	})
}
