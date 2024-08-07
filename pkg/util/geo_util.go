package util

import (
	"math"
)

// This is the formula to calculate the new latitude and longitude given a distance in kilometers and a bearing
// This util functions will mainly be used for testing purposes in order to easy the calculation of new coordinates
// and simplify the tests

// Earth radius in kilometers
const earthRadiusKm = 6371.0

// Add distance in kilometers to latitude and longitude
func AddKM(lat, lon, distanceKm, bearing float64) (float64, float64) {
	// degrees to rads
	latRad := degreesToRadians(lat)
	lonRad := degreesToRadians(lon)

	// bearing to rads
	bearingRad := degreesToRadians(bearing)

	// Convertir distancia to rad
	distanceRad := distanceKm / earthRadiusKm

	// Calculate new latitude
	newLatRad := math.Asin(math.Sin(latRad)*math.Cos(distanceRad) + math.Cos(latRad)*math.Sin(distanceRad)*math.Cos(bearingRad))

	// Calculate new longitude
	newLonRad := lonRad + math.Atan2(math.Sin(bearingRad)*math.Sin(distanceRad)*math.Cos(latRad), math.Cos(distanceRad)-math.Sin(latRad)*math.Sin(newLatRad))

	// Get back to degrees
	newLat := radiansToDegrees(newLatRad)
	newLon := radiansToDegrees(newLonRad)

	return newLat, newLon
}

// Convert from degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// Convert from radians to degrees
func radiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}
