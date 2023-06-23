package routing

import "math"

// :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
// :::                                                                         :::
// :::  This routine calculates the distance between two points (given the     :::
// :::  latitude/longitude of those points). It is based on free code used to  :::
// :::  calculate the distance between two locations using GeoDataSource(TM)   :::
// :::  products.                                                              :::
// :::                                                                         :::
// :::  Definitions:                                                           :::
// :::    South latitudes are negative, east longitudes are positive           :::
// :::                                                                         :::
// :::  Passed to function:                                                    :::
// :::    lat1, lon1 = Latitude and Longitude of point 1 (in decimal degrees)  :::
// :::    lat2, lon2 = Latitude and Longitude of point 2 (in decimal degrees)  :::
// :::    optional: unit = the unit you desire for results                     :::
// :::           where: 'M' is statute miles (default, or omitted)             :::
// :::                  'K' is kilometers                                      :::
// :::                  'N' is nautical miles                                  :::
// :::                                                                         :::
// :::  Worldwide cities and other features databases with latitude longitude  :::
// :::  are available at https://www.geodatasource.com                         :::
// :::                                                                         :::
// :::  For enquiries, please contact sales@geodatasource.com                  :::
// :::                                                                         :::
// :::  Official Web site: https://www.geodatasource.com                       :::
// :::                                                                         :::
// :::          Golang code James Robert Perih (c) All Rights Reserved 2018    :::
// :::                                                                         :::
// :::           GeoDataSource.com (C) All Rights Reserved 2017                :::
// :::                                                                         :::
// :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
func calculateGEODistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	radLat1 := math.Pi * lat1 / 180
	radLat2 := math.Pi * lat2 / 180

	theta := lng1 - lng2
	radTheta := math.Pi * theta / 180

	dist := math.Sin(radLat1)*math.Sin(radLat2) + math.Cos(radLat1)*math.Cos(radLat2)*math.Cos(radTheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	return dist * 1.609344
}
