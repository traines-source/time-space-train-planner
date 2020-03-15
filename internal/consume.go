package internal

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"traines.eu/time-space-train-planner/providers"
	"traines.eu/time-space-train-planner/providers/dbtimetables"
)

type Consumer struct {
	providers        []providers.Provider
	providerStations []providers.ProviderStation
	stations         map[int]*Station
	lines            map[int]*Line
}

func (c *Consumer) RequestStationDataUntil(station providers.ProviderStation) time.Time {
	d, _ := time.ParseDuration("1h")
	return time.Now().Add(d)
}

func (c *Consumer) Stations() []providers.ProviderStation {
	return c.providerStations
}

func (c *Consumer) StationByName(name string) (providers.ProviderStation, error) {
	for _, v := range c.providerStations {
		if name == v.Name {
			return v, nil
		}
	}
	return providers.ProviderStation{}, errors.New("not found")
}

func (c *Consumer) StationByEva(evaNumber int) (providers.ProviderStation, error) {
	for _, v := range c.providerStations {
		if evaNumber == v.EvaNumber {
			return v, nil
		}
	}
	return providers.ProviderStation{}, errors.New("not found")
}

func (c *Consumer) UpsertStation(e providers.ProviderStation) {
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

func (c *Consumer) UpsertLine(e providers.ProviderLine) {
	val, ok := c.lines[e.ID]
	if !ok {
		val = &Line{ID: e.ID, Stops: map[*Station]*LineStop{}}
		c.lines[e.ID] = val
	}
	val.Name = e.Name
	val.Type = e.Type
	val.Message = e.Message
}

func (c *Consumer) UpsertLineStop(e providers.ProviderLineStop) {
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
		copyStopInfo(e.Planned, &val.Planned)
	}
	if e.Planned != nil {
		copyProviderStopInfo(e.Planned, &val.Planned)
	}
	val.Message = e.Message
}

func copyProviderStopInfo(from *providers.ProviderLineStopInfo, to *StopInfo) {
	to.Departure = from.Departure
	to.Arrival = from.Arrival
	to.DepartureTrack = from.Track
}

func (c *Consumer) CallProviders() {
	c.providers = []providers.Provider{&dbtimetables.Timetables{}}
	c.providerStations = defaultStations(8000240, 8070004, 8070003)
	c.stations = map[int]*Station{}
	c.lines = map[int]*Line{}

	for _, p := range c.providers {
		p.Fetch(c)
	}
	fmt.Printf("%+v", c.stations)
	fmt.Printf("%+v", c.lines)
	c.generateEdges()
}

func defaultStations(evaNumbers ...int) []providers.ProviderStation {
	var stations []providers.ProviderStation
	for _, n := range evaNumbers {
		stations = append(stations, providers.ProviderStation{EvaNumber: n})
	}
	return stations
}

func (c *Consumer) generateEdges() {
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
		}
	}
	for _, station := range c.stations {
		sort.Slice(station.Departures, func(i, j int) bool {
			return station.Departures[i].Actual.Departure.Before(station.Departures[j].Actual.Departure)
		})
	}

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
