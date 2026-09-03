package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"my_vehicles_an_backup_agent/app/facades"
)

type M20260903000002CreateVehicleServiceTypesTable struct{}

// Signature The unique signature for the migration.
func (r *M20260903000002CreateVehicleServiceTypesTable) Signature() string {
	return "20260903000002_create_vehicle_service_types_table"
}

// Up Run the migrations.
func (r *M20260903000002CreateVehicleServiceTypesTable) Up() error {
	if facades.Schema().HasTable("vehicle_service_types") {
		return nil
	}

	return facades.Schema().Create("vehicle_service_types", func(table schema.Blueprint) {
		table.String("id", 36)
		table.String("icon")
		table.String("label")

		table.Primary("id")
	})
}

// Down Reverse the migrations.
func (r *M20260903000002CreateVehicleServiceTypesTable) Down() error {
	return facades.Schema().DropIfExists("vehicle_service_types")
}
