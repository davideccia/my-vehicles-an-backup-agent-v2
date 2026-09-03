package routes

import (
	"github.com/goravel/framework/contracts/route"

	"my_vehicles_an_backup_agent/app/facades"
	"my_vehicles_an_backup_agent/app/http/controllers"
	"my_vehicles_an_backup_agent/app/http/middleware"
)

// Api registers the backup agent's contract with the Android app: a public,
// version-gated /ping, and an API-key + version gated /dashboard/stats and
// /backup/{table} group. It also registers /healthcheck, a plain liveness
// probe for the Docker HEALTHCHECK directive — unlike /ping, it carries no
// client contract, so it skips both the API key and schema version gates.
func Api() {
	healthController := controllers.NewHealthController()
	facades.Route().Get("/healthcheck", healthController.Check)

	pingController := controllers.NewPingController()
	facades.Route().Middleware(middleware.SchemaVersion()).Get("/ping", pingController.Ping)

	facades.Route().
		Middleware(middleware.ApiKey(), middleware.SchemaVersion()).
		Group(func(router route.Router) {
			dashboardController := controllers.NewDashboardController()
			router.Get("/dashboard/stats", dashboardController.Stats)

			backupController := controllers.NewBackupController()
			router.Post("/backup/vehicles", backupController.PostVehicles)
			router.Get("/backup/vehicles", backupController.GetVehicles)
			router.Post("/backup/vehicle_refuels", backupController.PostVehicleRefuels)
			router.Get("/backup/vehicle_refuels", backupController.GetVehicleRefuels)
			router.Post("/backup/vehicle_service_types", backupController.PostVehicleServiceTypes)
			router.Get("/backup/vehicle_service_types", backupController.GetVehicleServiceTypes)
			router.Post("/backup/vehicle_services", backupController.PostVehicleServices)
			router.Get("/backup/vehicle_services", backupController.GetVehicleServices)
			router.Post("/backup/vehicle_service_reminders", backupController.PostVehicleServiceReminders)
			router.Get("/backup/vehicle_service_reminders", backupController.GetVehicleServiceReminders)
		})
}
