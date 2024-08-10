package billing

import (
	"github.com/OscarMoya/Glubber/pkg/model"
	"github.com/OscarMoya/Glubber/pkg/util"
)

type Biller interface {
	EstimateRide(ride *model.Ride) error
}

type SimpleBiller struct {
	baseCost float64
	kmCharge float64
}

func NewSimpleBiller(baseCost, kmCharge float64) *SimpleBiller {
	return &SimpleBiller{
		baseCost: baseCost,
		kmCharge: kmCharge,
	}
}

func (sb *SimpleBiller) EstimateRide(ride *model.Ride) error {
	distance := util.CalculateDistance(ride.SrcLat, ride.SrcLon, ride.DstLat, ride.DstLon)
	ride.Price = sb.baseCost + (distance * sb.kmCharge)
	return nil
}
