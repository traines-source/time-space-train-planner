package internal

import (
	"log"
	"math"
	"sort"
	"time"
)

func relevantDeparture(stop *LineStop) time.Time {
	return stop.Planned.Departure
}

func (c *consumer) generateTimetableEdges() {
	for _, line := range c.lines {
		var stops []*LineStop = line.Stops
		sort.Slice(stops, func(i, j int) bool {
			return relevantDeparture(stops[i]).Before(relevantDeparture(stops[j]))
		})
		var lastNonCancelled *LineStop
		for i := 1; i < len(stops); i++ {
			if !stops[i-1].Cancelled && stops[i].Cancelled {
				lastNonCancelled = stops[i-1]
			}
			if !stops[i].Cancelled && lastNonCancelled != nil {
				c.generateTimetableEdge(lastNonCancelled, stops[i], line)
				lastNonCancelled = nil
			}
			if geoDistStations(stops[i-1].Station, stops[i].Station) == 0 {
				continue
			}
			c.generateTimetableEdge(stops[i-1], stops[i], line)
		}
	}
}

func (c *consumer) generateTimetableEdge(a *LineStop, b *LineStop, line *Line) {
	edge := &Edge{
		Line:      line,
		From:      a.Station,
		To:        b.Station,
		Message:   a.Message, // TODO both msgs?
		Cancelled: a.Cancelled || b.Cancelled,
		Redundant: true, // set by stost
	}
	copyStopInfo(&a.Planned, &b.Planned, &edge.Planned)
	copyStopInfo(&a.Current, &b.Current, &edge.Current)
	copyStopInfo(&a.Planned, &b.Planned, &edge.Actual)
	copyStopInfo(&a.Current, &b.Current, &edge.Actual)

	edge.Current.DepartureTrack = a.Current.DepartureTrack
	edge.Current.ArrivalTrack = b.Current.ArrivalTrack
	edge.Current.Departure = a.Current.Departure
	edge.Current.Arrival = b.Current.Arrival

	line.Route = append(line.Route, edge)
	edge.From.Departures = append(edge.From.Departures, edge)
	edge.To.Arrivals = append(edge.To.Arrivals, edge)
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func geoDistStations(from *Station, to *Station) float64 {
	return geoDist(from.Lat, from.Lon, to.Lat, to.Lon)
}

func geoDist(fromLat float64, fromLon float64, toLat float64, toLon float64) float64 {
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

func (c *consumer) sortEdges() { // TODO zero minute edges order?
	for _, station := range c.stations {
		sort.Slice(station.Departures, func(i, j int) bool {
			return station.Departures[i].Actual.Departure.Before(station.Departures[j].Actual.Departure)
		})
		sort.Slice(station.Arrivals, func(i, j int) bool {
			return station.Arrivals[i].Actual.Arrival.Before(station.Arrivals[j].Actual.Arrival)
		})
	}
}

func (c *consumer) generateEdges(from *Station, to *Station) *ErrorCode {
	if from == nil || to == nil {
		log.Print("Error: Origin or destination does not exist or does not have any arrivals/departures")
		return &ErrorCode{Code: 400, Msg: "Origin or destination does not exist or does not have any arrivals/departures"}
	}
	c.generateTimetableEdges()
	c.sortEdges()
	return nil
}
