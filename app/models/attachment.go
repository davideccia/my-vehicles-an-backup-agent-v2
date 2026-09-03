package models

// Attachment mirrors the `attachments` Room table (my_vehicles_an/DB.md).
// It has no DB-level FK: owner_id references vehicles.id or vehicle_services.id
// polymorphically, disambiguated by owner_type (see AttachmentOwnerType* constants).
type Attachment struct {
	ID         string `gorm:"column:id;primaryKey;type:varchar(36)" json:"id"`
	OwnerType  string `gorm:"column:owner_type;not null" json:"owner_type"`
	OwnerID    string `gorm:"column:owner_id;not null" json:"owner_id"`
	Collection string `gorm:"column:collection;not null" json:"collection"` // e.g. 'vehicles:generic', 'vehicle_services:receipt'
	FileName   string `gorm:"column:file_name;not null" json:"file_name"`
	MimeType   string `gorm:"column:mime_type;not null" json:"mime_type"`
	SizeBytes  int64  `gorm:"column:size_bytes;not null" json:"size_bytes"`
	LocalPath  string `gorm:"column:local_path;not null" json:"local_path"`
	// CreatedAt is client-provided (epoch millis), not server-managed —
	// disable GORM's implicit autoCreateTime convention for this field name.
	CreatedAt int64 `gorm:"column:created_at;not null;autoCreateTime:false" json:"created_at"`
}

func (r *Attachment) TableName() string {
	return "attachments"
}
