package model

import "github.com/OscarMoya/Glubber/pkg/authentication"

type InputMessage struct {
	Payload []byte
}

type OutputMessage struct {
	Payload []byte
	IsError bool
}

type DriverInputMessage struct {
	InputMessage
	DriverAuth *authentication.DriverClaims
}

type DriverOutputMessage struct {
	OutputMessage
}
