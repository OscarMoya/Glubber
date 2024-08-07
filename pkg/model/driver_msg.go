package model

type DriverMsgType string

const (
	// DriverLocationMsg is the message type for driver location updates
	DriverLocationMsg DriverMsgType = "driver_location"
	// DriveRequestMsg is the message type for passenger ride requests
	DriveRequestMsg DriverMsgType = "drive_request"
)

type (

	// BaseMessage is the base structure for all messages with a Type field
	BaseMessage struct {
		Type DriverMsgType `json:"type"`
	}

	// DriverLocation represents the location message from a driver
	DriverLocation struct {
		BaseMessage
		DriverID  string  `json:"driver_id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	// DriveRequest represents a ride request message from a passenger
	DriveRequest struct {
		BaseMessage
		PassengerID string  `json:"passenger_id"`
		PickupLat   float64 `json:"pickup_latitude"`
		PickupLng   float64 `json:"pickup_longitude"`
		DropLat     float64 `json:"drop_latitude"`
		DropLng     float64 `json:"drop_longitude"`
	}
)
