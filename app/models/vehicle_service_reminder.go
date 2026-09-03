package models

// VehicleServiceReminder mirrors the `vehicle_service_reminders` Room table (my_vehicles_an/DB.md).
type VehicleServiceReminder struct {
	ID                   string `gorm:"column:id;primaryKey;type:varchar(36)" json:"id"`
	VehicleID            string `gorm:"column:vehicle_id;not null" json:"vehicle_id"`
	VehicleServiceTypeID string `gorm:"column:vehicle_service_type_id;not null" json:"vehicle_service_type_id"`
	Every                int    `gorm:"column:every;not null" json:"every"` // interval in km or days
}

func (r *VehicleServiceReminder) TableName() string {
	return "vehicle_service_reminders"
}
