package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/OscarMoya/Glubber/pkg/location"
	"github.com/OscarMoya/Glubber/pkg/model"
)

func driverSvcLoop(ctx context.Context, in <-chan *model.DriverInputMessage, out chan<- *model.DriverOutputMessage, geoService location.LocationManager) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Driver service loop done")
			return
		case msg := <-in:
			// TODO: Stablish timeout for this OPS
			backendCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			var baseMessage model.BaseMessage
			if err := json.Unmarshal(msg.Payload, &baseMessage); err != nil {
				log.Println("unmarshal base message:", err)
				continue
			}

			switch baseMessage.Type {
			case model.DriverLocationMsgType:
				var loc model.DriverLocationRequest
				if err := json.Unmarshal(msg.Payload, &loc); err != nil {
					log.Println("unmarshal driver location:", err)
					continue
				}
				go handleDriverLocation(backendCtx, out, msg.DriverAuth.DriverID, loc.Latitude, loc.Longitude, geoService)

			case model.DriveRequestMsgType:
				var req model.DriveRequest
				if err := json.Unmarshal(msg.Payload, &req); err != nil {
					log.Println("unmarshal drive request:", err)
					continue
				}
				handleDriveRequest(backendCtx, out, req)

			default:
				log.Println("Unknown message type:", baseMessage.Type)
			}
		}
	}
}

func handleDriverLocation(ctx context.Context, out chan<- *model.DriverOutputMessage, driverID string, latitude, longitude float64, geoService location.LocationManager) {
	err := geoService.SaveDriverLocation(ctx, driverID, latitude, longitude)
	if err != nil {
		log.Println("SaveDriverLocation:", err)
		errMsg := model.DriverErrorResponse{
			Code:   500,
			Reason: err.Error(),
		}
		errMsg.Type = model.DriverErrorResponseMsgType
		payload, err := json.Marshal(errMsg)
		if err != nil {
			log.Println("marshal error response:", err)
			return
		}
		outMsg := &model.DriverOutputMessage{}
		outMsg.IsError = true
		outMsg.Payload = payload

		out <- outMsg

		return
	}
}

func handleDriveRequest(ctx context.Context, out chan<- *model.DriverOutputMessage, req model.DriveRequest) {
	log.Printf("Received drive request: %+v\n", req)
}
