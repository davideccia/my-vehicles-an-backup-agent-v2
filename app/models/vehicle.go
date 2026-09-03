package models

// Vehicle mirrors the `vehicles` Room table (my_vehicles_an/DB.md).
type Vehicle struct {
	ID               string `gorm:"column:id;primaryKey;type:varchar(36)" json:"id"`
	PlateNumber      string `gorm:"column:plate_number;not null;unique" json:"plate_number"`
	Brand            string `gorm:"column:brand;not null" json:"brand"`
	Model            string `gorm:"column:model;not null" json:"model"`
	Year             int    `gorm:"column:year;not null" json:"year"`
	FuelType         string `gorm:"column:fuel_type;not null" json:"fuel_type"` // PETROL | DIESEL
	FuelUnit         string `gorm:"column:fuel_unit;not null" json:"fuel_unit"` // LITERS
	PurchaseDate     *int64 `gorm:"column:purchase_date" json:"purchase_date"`  // LocalDate as epoch-day
	PurchaseOdometer *int64 `gorm:"column:purchase_odometer" json:"purchase_odometer"`
	Color            string `gorm:"column:color;not null;default:'#4CAF50'" json:"color"`
}

func (r *Vehicle) TableName() string {
	return "vehicles"
}
