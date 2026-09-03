package models

// VehicleRefuel mirrors the `vehicle_refuels` Room table (my_vehicles_an/DB.md).
type VehicleRefuel struct {
	ID         string  `gorm:"column:id;primaryKey;type:varchar(36)" json:"id"`
	VehicleID  string  `gorm:"column:vehicle_id;not null" json:"vehicle_id"`
	Date       int64   `gorm:"column:date;not null" json:"date"` // LocalDate as epoch-day
	UnitPrice  float64 `gorm:"column:unit_price;not null" json:"unit_price"`
	Liters     float64 `gorm:"column:liters;not null" json:"liters"`
	TotalPrice float64 `gorm:"column:total_price;not null" json:"total_price"`
	TankStatus *int    `gorm:"column:tank_status" json:"tank_status"` // 0-100%, nullable
	Place      *string `gorm:"column:place" json:"place"`
	Odometer   *int64  `gorm:"column:odometer" json:"odometer"`
}

func (r *VehicleRefuel) TableName() string {
	return "vehicle_refuels"
}
