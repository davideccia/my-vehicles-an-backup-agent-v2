package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"my_vehicles_an_backup_agent/app/facades"
)

type M20260903000004CreateVehicleServicesTable struct{}

// Signature The unique signature for the migration.
func (r *M20260903000004CreateVehicleServicesTable) Signature() string {
	return "20260903000004_create_vehicle_services_table"
}

// Up Run the migrations.
func (r *M20260903000004CreateVehicleServicesTable) Up() error {
	if facades.Schema().HasTable("vehicle_services") {
		return nil
	}

	return facades.Schema().Create("vehicle_services", func(table schema.Blueprint) {
		table.String("id", 36)
		table.String("vehicle_id")
		table.String("vehicle_service_type_id")
		table.BigInteger("date") // LocalDate as epoch-day
		table.Double("amount")
		table.String("place").Nullable()
		table.String("notes").Nullable()
		table.BigInteger("odometer").Nullable()

		table.Primary("id")
		table.Index("vehicle_id")
		table.Index("vehicle_service_type_id")
		table.Index("date")
	})
}

// Down Reverse the migrations.
func (r *M20260903000004CreateVehicleServicesTable) Down() error {
	return facades.Schema().DropIfExists("vehicle_services")
}
