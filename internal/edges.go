package internal

import (
	"fmt"
	"math"
	"sort"
	"time"
)

const maxFootDistMeters = 5000
const footKmh = 6

func (c *consumer) generateTimetableEdges() {
	for _, line := range c.lines {
		var stops []*LineStop = line.Stops
		sort.Slice(stops, func(i, j int) bool {
			if !stops[i].Current.Departure.IsZero() && !stops[j].Current.Departure.IsZero() {
				return stops[i].Current.Departure.Before(stops[j].Current.Departure)
			}
			return stops[i].Planned.Departure.Before(stops[j].Planned.Departure)
		})
		for i := 1; i < len(stops); i++ {
			if geoDistStations(stops[i-1].Station, stops[i].Station) == 0 {
				continue
			}
			edge := &Edge{
				Line:    line,
				From:    stops[i-1].Station,
				To:      stops[i].Station,
				Message: stops[i-1].Message,
			}
			copyStopInfo(&stops[i-1].Planned, &stops[i].Planned, &edge.Planned)
			copyStopInfo(&stops[i-1].Current, &stops[i].Current, &edge.Current)
			copyStopInfo(&stops[i-1].Planned, &stops[i].Planned, &edge.Actual)
			copyStopInfo(&stops[i-1].Current, &stops[i].Current, &edge.Actual)

			edge.Current.DepartureTrack = stops[i-1].Current.DepartureTrack
			edge.Current.ArrivalTrack = stops[i].Current.ArrivalTrack
			edge.Current.Departure = stops[i-1].Current.Departure
			edge.Current.Arrival = stops[i].Current.Arrival

			line.Route = append(line.Route, edge)
			edge.From.Departures = append(edge.From.Departures, edge)
			edge.To.Arrivals = append(edge.To.Arrivals, edge)
		}
	}
}

func (c *consumer) generateOnFootEdges(origin *Station, destination *Station) {
	for _, s1 := range c.stations {
		for _, s2 := range c.stations {
			if s1 == s2 {
				continue
			}
			var dist = geoDistStations(s1, s2)
			if dist > maxFootDistMeters {
				continue
			}
			c.generateOnFootEdgesBetweenTwoStationsInDirection(s1, s2, dist, origin, destination)
		}
	}
}

func (c *consumer) generateOnFootEdgesBetweenTwoStationsInDirection(from *Station, to *Station, dist float64, origin *Station, destination *Station) {
	if from == destination || to == origin {
		return
	}
	correspondances := to.Departures
	if to == destination {
		correspondances = from.Arrivals
	}
	var duration = time.Minute * time.Duration(math.Ceil(dist/1000/footKmh*60))
	for _, correspondance := range correspondances {
		if correspondance.Line.Type == "Foot" {
			continue
		}
		var departure time.Time
		var arrival time.Time
		if to == destination {
			departure = correspondance.Actual.Arrival
			arrival = correspondance.Actual.Arrival.Add(duration)
		} else {
			departure = correspondance.Actual.Departure.Add(-duration)
			arrival = correspondance.Actual.Departure
			if len(from.Arrivals) == 0 || departure.Before(from.Arrivals[0].Actual.Arrival) {
				continue
			}
		}
		c.generateOnFootEdgeBetweenTwoStationsInDirection(from, to, dist, departure, arrival)
	}
}

func (c *consumer) generateOnFootEdgeBetweenTwoStationsInDirection(from *Station, to *Station, dist float64, departure time.Time, arrival time.Time) {

	var lineID = fmt.Sprint(int64(from.EvaNumber*1000000000+to.EvaNumber*100) + departure.Unix())

	var line = &Line{
		ID:   lineID,
		Name: fmt.Sprintf("%.0fm", math.Round(dist)),
		Type: "Foot",
	}
	c.lines[lineID] = line
	edge := &Edge{
		Line: line,
		From: from,
		To:   to,
		Actual: StopInfo{
			Departure: departure,
			Arrival:   arrival,
		},
	}
	line.Route = append(line.Route, edge)
	edge.From.Departures = append(edge.From.Departures, edge)
	edge.To.Arrivals = append(edge.To.Arrivals, edge)
}

func degreesToRadians(degrees float32) float64 {
	return float64(degrees) * math.Pi / 180
}

func geoDistStations(from *Station, to *Station) float64 {
	return geoDist(from.Lat, from.Lon, to.Lat, to.Lon)
}

func geoDist(fromLat float32, fromLon float32, toLat float32, toLon float32) float64 {
	var earthRadiusM float64 = 6371000

	var dLat = degreesToRadians(toLat - fromLat)
	var dLon = degreesToRadians(toLon - fromLon)

	var fromLatRad = degreesToRadians(fromLat)
	var toLatRad = degreesToRadians(toLat)

	var a = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(fromLatRad)*math.Cos(toLatRad)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusM * c
}

func (c *consumer) sortEdges() {
	for _, station := range c.stations {
		sort.Slice(station.Departures, func(i, j int) bool {
			return station.Departures[i].Actual.Departure.Before(station.Departures[j].Actual.Departure)
		})
		sort.Slice(station.Arrivals, func(i, j int) bool {
			return station.Arrivals[i].Actual.Arrival.Before(station.Arrivals[j].Actual.Arrival)
		})
	}
}

func (c *consumer) generateEdges(from *Station, to *Station) {
	c.generateTimetableEdges()
	c.sortEdges()
	c.generateOnFootEdges(from, to)
	c.sortEdges()
}
