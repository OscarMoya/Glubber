package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/OscarMoya/Glubber/pkg/location"
	"github.com/OscarMoya/Glubber/pkg/model"
	"github.com/gorilla/websocket"
)

func handleDriverConnections(w http.ResponseWriter, r *http.Request, serviceStatus *ServiceData) {

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	authenticator := serviceStatus.Authenticator
	claims, err := authenticator.ValidateDriverJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	in := make(chan *model.DriverInputMessage, 256)
	out := make(chan *model.DriverOutputMessage, 256)

	defer ws.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go driverSvcLoop(ctx, in, out, serviceStatus.GeoService)

	for {
		select {
		case <-ctx.Done():
			log.Println("Driver service loop done")
			return
		case msg := <-out:
			err := ws.WriteMessage(websocket.TextMessage, msg.Payload)
			if err != nil {
				// If there is an error writing to the websocket, log it and continue
				// We may want to handle this differently in a production system
				log.Println("write:", err)
				continue
			}
		default:
			// TODO: check message types
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			// TODO: Add constructor
			driverIn := &model.DriverInputMessage{}
			driverIn.Payload = msg
			driverIn.DriverAuth = claims

			in <- driverIn
		}
	}
}

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

			case model.DriverGoodByeMsgType:
				go handleDriverGoodBye(backendCtx, out, msg.DriverAuth.DriverID, geoService)

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

func handleDriverGoodBye(ctx context.Context, out chan<- *model.DriverOutputMessage, driverID string, geoService location.LocationManager) {
	err := geoService.RemoveDriverLocation(ctx, driverID)
	if err != nil {
		log.Println("DeleteDriverLocation:", err)
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
