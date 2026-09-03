package models

// The types below are read-only DAO projections (my_vehicles_an/DB.md), not
// persisted tables — they have no TableName() and are meant as scan targets
// for aggregate queries, e.g.:
//   facades.Orm().Query().Model(&VehicleRefuel{}).
//     Select("vehicle_id", "MAX(odometer) as max_odometer").
//     Group("vehicle_id").Scan(&results)

// VehicleMaxOdometer is the highest recorded odometer reading for a vehicle.
type VehicleMaxOdometer struct {
	VehicleID   string `json:"vehicle_id"`
	MaxOdometer int64  `json:"max_odometer"`
}

// ServiceTypeLastOdometer is the odometer reading at the most recent service
// of a given type for a vehicle.
type ServiceTypeLastOdometer struct {
	VehicleID            string `json:"vehicle_id"`
	VehicleServiceTypeID string `json:"vehicle_service_type_id"`
	LastOdometer         int64  `json:"last_odometer"`
}

// ServiceTypeLastService is the most recent service of a given type for a
// vehicle, with its odometer, date and amount.
type ServiceTypeLastService struct {
	VehicleID            string   `json:"vehicle_id"`
	VehicleServiceTypeID string   `json:"vehicle_service_type_id"`
	LastOdometer         int64    `json:"last_odometer"`
	LastDate             *int64   `json:"last_date"`
	LastAmount           *float64 `json:"last_amount"`
}
