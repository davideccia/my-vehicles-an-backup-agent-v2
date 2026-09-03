package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"my_vehicles_an_backup_agent/app/facades"
)

type M20260903000003CreateVehicleRefuelsTable struct{}

// Signature The unique signature for the migration.
func (r *M20260903000003CreateVehicleRefuelsTable) Signature() string {
	return "20260903000003_create_vehicle_refuels_table"
}

// Up Run the migrations.
func (r *M20260903000003CreateVehicleRefuelsTable) Up() error {
	if facades.Schema().HasTable("vehicle_refuels") {
		return nil
	}

	return facades.Schema().Create("vehicle_refuels", func(table schema.Blueprint) {
		table.String("id", 36)
		table.String("vehicle_id")
		table.BigInteger("date") // LocalDate as epoch-day
		table.Double("unit_price")
		table.Double("liters")
		table.Double("total_price")
		table.Integer("tank_status").Nullable() // 0-100%
		table.String("place").Nullable()
		table.BigInteger("odometer").Nullable()

		table.Primary("id")
		table.Index("vehicle_id")
		table.Index("date")
	})
}

// Down Reverse the migrations.
func (r *M20260903000003CreateVehicleRefuelsTable) Down() error {
	return facades.Schema().DropIfExists("vehicle_refuels")
}
