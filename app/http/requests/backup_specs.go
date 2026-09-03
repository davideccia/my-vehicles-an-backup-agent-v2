package requests

import "my_vehicles_an_backup_agent/app/models"

// The specs below mirror the class-validator decorators on the NestJS
// reference server's DTOs (src/backup/dto/backup.dto.ts), field by field.

var VehicleSpec = TableSpec{Fields: []FieldSpec{
	{Name: "id", Required: true},
	{Name: "plate_number", Required: true},
	{Name: "brand", Required: true},
	{Name: "model", Required: true},
	{Name: "year", Required: true},
	{Name: "fuel_type", Required: true, Enum: []string{models.FuelTypePetrol, models.FuelTypeDiesel}},
	{Name: "fuel_unit", Required: true, Enum: []string{models.FuelUnitLiters}},
	{Name: "purchase_date"},
	{Name: "purchase_odometer"},
	{Name: "color", Required: true},
}}

var VehicleRefuelSpec = TableSpec{Fields: []FieldSpec{
	{Name: "id", Required: true},
	{Name: "vehicle_id", Required: true},
	{Name: "date", Required: true},
	{Name: "unit_price", Required: true},
	{Name: "liters", Required: true},
	{Name: "total_price", Required: true},
	{Name: "tank_status"},
	{Name: "place"},
	{Name: "odometer"},
}}

var VehicleServiceTypeSpec = TableSpec{Fields: []FieldSpec{
	{Name: "id", Required: true},
	{Name: "icon", Required: true},
	{Name: "label", Required: true},
}}

var VehicleServiceSpec = TableSpec{Fields: []FieldSpec{
	{Name: "id", Required: true},
	{Name: "vehicle_id", Required: true},
	{Name: "vehicle_service_type_id", Required: true},
	{Name: "date", Required: true},
	{Name: "amount", Required: true},
	{Name: "place"},
	{Name: "notes"},
	{Name: "odometer"},
}}

var VehicleServiceReminderSpec = TableSpec{Fields: []FieldSpec{
	{Name: "id", Required: true},
	{Name: "vehicle_id", Required: true},
	{Name: "vehicle_service_type_id", Required: true},
	{Name: "every", Required: true},
}}
