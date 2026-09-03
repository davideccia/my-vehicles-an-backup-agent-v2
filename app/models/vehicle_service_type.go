package models

// VehicleServiceType mirrors the `vehicle_service_types` Room table (my_vehicles_an/DB.md).
type VehicleServiceType struct {
	ID    string `gorm:"column:id;primaryKey;type:varchar(36)" json:"id"`
	Icon  string `gorm:"column:icon;not null" json:"icon"`
	Label string `gorm:"column:label;not null" json:"label"`
}

func (r *VehicleServiceType) TableName() string {
	return "vehicle_service_types"
}
