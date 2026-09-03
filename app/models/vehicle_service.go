package models

// VehicleService mirrors the `vehicle_services` Room table (my_vehicles_an/DB.md).
type VehicleService struct {
	ID                   string  `gorm:"column:id;primaryKey;type:varchar(36)" json:"id"`
	VehicleID            string  `gorm:"column:vehicle_id;not null" json:"vehicle_id"`
	VehicleServiceTypeID string  `gorm:"column:vehicle_service_type_id;not null" json:"vehicle_service_type_id"`
	Date                 int64   `gorm:"column:date;not null" json:"date"` // LocalDate as epoch-day
	Amount               float64 `gorm:"column:amount;not null" json:"amount"`
	Place                *string `gorm:"column:place" json:"place"`
	Notes                *string `gorm:"column:notes" json:"notes"`
	Odometer             *int64  `gorm:"column:odometer" json:"odometer"`
}

func (r *VehicleService) TableName() string {
	return "vehicle_services"
}
