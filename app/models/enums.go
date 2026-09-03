package models

// Fuel type/unit values, mirrored from the FuelType/FuelUnit converters
// in the Android app (my_vehicles_an), stored as TEXT.
const (
	FuelTypePetrol = "PETROL"
	FuelTypeDiesel = "DIESEL"

	FuelUnitLiters = "LITERS"
)

// Attachment owner types, mirrored from AttachmentOwnerType.kt in the
// Android app (my_vehicles_an/app/src/main/java/com/davideccia/myvehicles/data/AttachmentOwnerType.kt).
const (
	AttachmentOwnerTypeVehicle        = "VEHICLE"
	AttachmentOwnerTypeVehicleService = "VEHICLE_SERVICE"
)

// Attachment collection keys, mirrored from Vehicle.kt and VehicleService.kt
// in the Android app.
const (
	AttachmentCollectionVehicleGeneric        = "vehicles:generic"
	AttachmentCollectionVehicleServiceReceipt = "vehicle_services:receipt"
	AttachmentCollectionVehicleServiceGeneric = "vehicle_services:attachments"
)
