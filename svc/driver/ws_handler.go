package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/OscarMoya/Glubber/pkg/model"
)

func driverSvcLoop(ctx context.Context, in <-chan []byte, out chan<- []byte) {
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
			if err := json.Unmarshal(msg, &baseMessage); err != nil {
				log.Println("unmarshal base message:", err)
				continue
			}

			switch baseMessage.Type {
			case model.DriverLocationMsg:
				var loc model.DriverLocation
				if err := json.Unmarshal(msg, &loc); err != nil {
					log.Println("unmarshal driver location:", err)
					continue
				}
				handleDriverLocation(backendCtx, out, loc)

			case model.DriveRequestMsg:
				var req model.DriveRequest
				if err := json.Unmarshal(msg, &req); err != nil {
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

func handleDriverLocation(ctx context.Context, out chan<- []byte, loc model.DriverLocation) {
	log.Printf("Received driver location: %+v\n", loc)
}

func handleDriveRequest(ctx context.Context, out chan<- []byte, req model.DriveRequest) {
	log.Printf("Received drive request: %+v\n", req)
}
