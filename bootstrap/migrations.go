package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"my_vehicles_an_backup_agent/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20260903000001CreateVehiclesTable{},
		&migrations.M20260903000002CreateVehicleServiceTypesTable{},
		&migrations.M20260903000003CreateVehicleRefuelsTable{},
		&migrations.M20260903000004CreateVehicleServicesTable{},
		&migrations.M20260903000005CreateVehicleServiceRemindersTable{},
		&migrations.M20260903000006CreateAttachmentsTable{},
	}
}
