package model

type DriverMsgType string

const (
	// DriverLocationMsgType is the message type for driver location updates
	DriverLocationMsgType DriverMsgType = "driver_location"
	// DriveRequestMsgType is the message type for passenger ride requests
	DriveRequestMsgType DriverMsgType = "driver_request"
	// DriverErrorResponseMsgType is the message type for error responses
	DriverErrorResponseMsgType DriverMsgType = "driver_error"
)

type (

	// BaseMessage is the base structure for all messages with a Type field
	BaseMessage struct {
		Type DriverMsgType `json:"type"`
	}

	// DriverLocation represents the location message from a driver
	// This message is sent from the Client to the Server
	DriverLocationRequest struct {
		BaseMessage
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	// DriveRequest represents a ride request message from a passenger
	// This message is sent from the Client to the Server
	DriveRequest struct {
		BaseMessage
		PickupLat float64 `json:"pickup_latitude"`
		PickupLng float64 `json:"pickup_longitude"`
		DropLat   float64 `json:"drop_latitude"`
		DropLng   float64 `json:"drop_longitude"`
	}

	// DriverHelloRequest represents a driver hello message
	// This message is sent from the Client to the Server
	DriverHelloRequest struct {
		BaseMessage
	}

	// DriverHelloResponse represents a driver hello response message
	// This message is sent from the Server to the Client
	DriverHelloResponse struct {
		BaseMessage
	}

	// DriverGoodByeRequest represents a driver goodbye message
	// This message is sent from the Client to the Server
	DriverGoodByeRequest struct {
		BaseMessage
	}

	// DriverGoodByeResponse represents a driver goodbye response message
	// This message is sent from the Server to the Client
	DriverGoodByeResponse struct {
		BaseMessage
	}

	// DriverErrorResponse represents an error response message
	// This message is sent from the Server to the Client
	DriverErrorResponse struct {
		BaseMessage
		OriginalMessageType DriverMsgType `json:"original_message_type"`
		Code                int           `json:"code"`
		Reason              string        `json:"reason"`
	}
)
