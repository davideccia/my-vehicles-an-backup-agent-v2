package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"my_vehicles_an_backup_agent/app/facades"
)

type M20260903000006CreateAttachmentsTable struct{}

// Signature The unique signature for the migration.
func (r *M20260903000006CreateAttachmentsTable) Signature() string {
	return "20260903000006_create_attachments_table"
}

// Up Run the migrations.
func (r *M20260903000006CreateAttachmentsTable) Up() error {
	if facades.Schema().HasTable("attachments") {
		return nil
	}

	return facades.Schema().Create("attachments", func(table schema.Blueprint) {
		table.String("id", 36)
		table.String("owner_type")
		table.String("owner_id")
		table.String("collection") // e.g. 'vehicles:generic', 'vehicle_services:receipt'
		table.String("file_name")
		table.String("mime_type")
		table.BigInteger("size_bytes")
		table.String("local_path")
		table.BigInteger("created_at") // epoch millis, client-provided

		table.Primary("id")
		table.Index("owner_type", "owner_id", "collection").Name("idx_attachments_owner")
	})
}

// Down Reverse the migrations.
func (r *M20260903000006CreateAttachmentsTable) Down() error {
	return facades.Schema().DropIfExists("attachments")
}
