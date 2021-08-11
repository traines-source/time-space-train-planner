package internal

import (
	"errors"
	"log"
	"sort"
	"time"

	"traines.eu/time-space-train-planner/providers"
	"traines.eu/time-space-train-planner/providers/dbtimetables"
)

type consumer struct {
	providers        []providers.Provider
	providerStations []providers.ProviderStation
	stations         map[int]*Station
	lines            map[int]*Line
}

func (c *consumer) RequestStationDataBetween(station *providers.ProviderStation) (from time.Time, to time.Time) {
	delta, _ := time.ParseDuration("3h")
	t := time.Now()
	from = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.Local)
	from = time.Date(t.Year(), t.Month(), 11, 19, 0, 0, 0, time.Local)
	return from, from.Add(delta)
}

func (c *consumer) Stations() []providers.ProviderStation {
	return c.providerStations
}

func (c *consumer) StationByName(name string) (providers.ProviderStation, error) {
	for _, v := range c.providerStations {
		if name == v.Name {
			return v, nil
		}
	}
	return providers.ProviderStation{}, errors.New("not found")
}

func (c *consumer) StationByEva(evaNumber int) (providers.ProviderStation, error) {
	for _, v := range c.providerStations {
		if evaNumber == v.EvaNumber {
			return v, nil
		}
	}
	return providers.ProviderStation{}, errors.New("not found")
}

func (c *consumer) UpsertStation(e providers.ProviderStation) {
	var station *providers.ProviderStation
	for _, s := range c.providerStations {
		if s.EvaNumber == e.EvaNumber {
			station = &s
		}
	}
	if station == nil {
		station = &providers.ProviderStation{EvaNumber: e.EvaNumber}
	}
	station.Name = e.Name
	station.Lat = e.Lat
	station.Lon = e.Lon

	val, ok := c.stations[e.EvaNumber]
	if !ok {
		val = &Station{EvaNumber: e.EvaNumber}
		c.stations[e.EvaNumber] = val
	}
	val.Name = e.Name
	val.Lat = e.Lat
	val.Lon = e.Lon
}

func (c *consumer) UpsertLine(e providers.ProviderLine) {
	val, ok := c.lines[e.ID]
	if !ok {
		val = &Line{ID: e.ID, Stops: map[*Station]*LineStop{}}
		c.lines[e.ID] = val
	}
	val.Name = e.Name
	val.Type = e.Type
	val.Message = e.Message
}

func (c *consumer) UpsertLineStop(e providers.ProviderLineStop) {
	station, ok := c.stations[e.EvaNumber]
	if !ok {
		log.Panicf("Non-existant Station %d for stop of Line %d", e.EvaNumber, e.LineID)
		return
	}
	line, ok := c.lines[e.LineID]
	if !ok {
		log.Panicf("Non-existant Line %d for Station  %d", e.LineID, e.EvaNumber)
		return
	}
	val, ok := line.Stops[station]
	if !ok {
		val = &LineStop{Station: station}
		line.Stops[station] = val
	}
	if e.Planned != nil {
		copyProviderStopInfo(e.Planned, &val.Planned)
	}
	if e.Current != nil {
		copyProviderStopInfo(e.Current, &val.Current)
	}
	val.Message = e.Message
}

func copyProviderStopInfo(from *providers.ProviderLineStopInfo, to *StopInfo) {
	to.Departure = from.Departure
	to.Arrival = from.Arrival
	if from.Departure.IsZero() {
		to.Departure = from.Arrival
	}
	if from.Arrival.IsZero() {
		to.Arrival = from.Departure
	}
	to.DepartureTrack = from.Track
}

func (c *consumer) callProviders() {
	c.providers = []providers.Provider{&dbtimetables.Timetables{}}
	c.providerStations = defaultStations(8003819, 8003816, 8000240, 8070004, 8070003, 8000257, 8000236, 8000244, 8000096)
	// 8000105, 8098105
	c.stations = map[int]*Station{}
	c.lines = map[int]*Line{}

	for _, p := range c.providers {
		p.Fetch(c)
	}
}

func defaultStations(evaNumbers ...int) []providers.ProviderStation {
	var stations []providers.ProviderStation
	for _, n := range evaNumbers {
		stations = append(stations, providers.ProviderStation{EvaNumber: n})
	}
	return stations
}

func (c *consumer) generateEdges() {
	for _, line := range c.lines {
		var stops []*LineStop
		for _, stop := range line.Stops {
			stops = append(stops, stop)
		}
		sort.Slice(stops, func(i, j int) bool {
			// TODO current?
			return stops[i].Planned.Departure.Before(stops[j].Planned.Departure)
		})
		for i := 1; i < len(stops); i++ {
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
			edge.Current.Departure = stops[i-1].Current.Departure
			edge.Current.Arrival = stops[i].Current.Arrival

			line.Route = append(line.Route, edge)
			edge.From.Departures = append(edge.From.Departures, edge)
			edge.To.Arrivals = append(edge.To.Arrivals, edge)
		}
	}
	for _, station := range c.stations {
		sort.Slice(station.Departures, func(i, j int) bool {
			return station.Departures[i].Actual.Departure.Before(station.Departures[j].Actual.Departure)
		})
		sort.Slice(station.Arrivals, func(i, j int) bool {
			return station.Arrivals[i].Actual.Departure.Before(station.Arrivals[j].Actual.Departure)
		})
	}

}

func (c *consumer) rankStations() (destination *Station) {
	i := 0
	for _, s := range c.providerStations {
		c.stations[s.EvaNumber].Rank = i
		destination = c.stations[s.EvaNumber]
		i++
	}
	return destination
}

func copyStopInfo(lastFrom *StopInfo, thisFrom *StopInfo, to *StopInfo) {
	if lastFrom.DepartureTrack != "" {
		to.DepartureTrack = lastFrom.DepartureTrack
	}
	if !lastFrom.Departure.IsZero() {
		to.Departure = lastFrom.Departure
	}
	if !thisFrom.Arrival.IsZero() {
		to.Arrival = thisFrom.Arrival
	}
}

func ObtainData() (map[int]*Station, map[int]*Line) {
	c := &consumer{}
	c.callProviders()
	c.generateEdges()
	destination := c.rankStations()
	shortestPaths(c.stations, destination)
	return c.stations, c.lines
}
