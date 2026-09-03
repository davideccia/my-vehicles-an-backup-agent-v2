package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"my_vehicles_an_backup_agent/app/facades"
)

type M20260903000001CreateVehiclesTable struct{}

// Signature The unique signature for the migration.
func (r *M20260903000001CreateVehiclesTable) Signature() string {
	return "20260903000001_create_vehicles_table"
}

// Up Run the migrations.
func (r *M20260903000001CreateVehiclesTable) Up() error {
	if facades.Schema().HasTable("vehicles") {
		return nil
	}

	return facades.Schema().Create("vehicles", func(table schema.Blueprint) {
		table.String("id", 36)
		table.String("plate_number")
		table.String("brand")
		table.String("model")
		table.Integer("year")
		table.String("fuel_type")
		table.String("fuel_unit")
		table.BigInteger("purchase_date").Nullable()
		table.BigInteger("purchase_odometer").Nullable()
		table.String("color").Default("#4CAF50")

		table.Primary("id")
		table.Unique("plate_number")
	})
}

// Down Reverse the migrations.
func (r *M20260903000001CreateVehiclesTable) Down() error {
	return facades.Schema().DropIfExists("vehicles")
}
