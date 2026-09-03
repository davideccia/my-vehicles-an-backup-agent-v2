package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"my_vehicles_an_backup_agent/app/facades"
)

type M20260903000005CreateVehicleServiceRemindersTable struct{}

// Signature The unique signature for the migration.
func (r *M20260903000005CreateVehicleServiceRemindersTable) Signature() string {
	return "20260903000005_create_vehicle_service_reminders_table"
}

// Up Run the migrations.
func (r *M20260903000005CreateVehicleServiceRemindersTable) Up() error {
	if facades.Schema().HasTable("vehicle_service_reminders") {
		return nil
	}

	return facades.Schema().Create("vehicle_service_reminders", func(table schema.Blueprint) {
		table.String("id", 36)
		table.String("vehicle_id")
		table.String("vehicle_service_type_id")
		table.Integer("every") // interval in km or days

		table.Primary("id")
		table.Index("vehicle_id")
		table.Index("vehicle_service_type_id")
	})
}

// Down Reverse the migrations.
func (r *M20260903000005CreateVehicleServiceRemindersTable) Down() error {
	return facades.Schema().DropIfExists("vehicle_service_reminders")
}
