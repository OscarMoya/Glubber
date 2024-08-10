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

// CalculateDistance calculates the distance between two points specified by their latitude and longitude in degrees using haversine formula.
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Radius of the Earth in kilometers

	// Convert degrees to radians
	lat1 = lat1 * math.Pi / 180
	lon1 = lon1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180
	lon2 = lon2 * math.Pi / 180

	// Differences
	dlat := lat2 - lat1
	dlon := lon2 - lon1

	// Haversine formula
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in kilometers
	distance := R * c
	return distance
}

// Convert from degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// Convert from radians to degrees
func radiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}
