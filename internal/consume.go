package internal

import (
	"errors"
	"log"
	"sort"
	"time"

	"traines.eu/time-space-train-planner/providers"
	"traines.eu/time-space-train-planner/providers/dbrest"
)

type consumer struct {
	providers        []providers.Provider
	providerStations []providers.ProviderStation
	stations         map[int]*Station
	lines            map[int]*Line
	dateTime         time.Time
}

func (c *consumer) RequestStationDataBetween(station *providers.ProviderStation) (from time.Time, to time.Time) {
	delta, _ := time.ParseDuration("4h")

	log.Print("Requesting for ", c.dateTime)
	//t := time.Now()
	//from = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.Local)
	//from = time.Date(t.Year(), t.Month(), 9, 19, 0, 0, 0, time.Local)
	from = c.dateTime
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
		//c.providerStations = append(c.providerStations, *station)
	}
	if e.Name != "" {
		station.Name = e.Name
	}
	if e.Lat != 0 {
		station.Lat = e.Lat
	}
	if e.Lon != 0 {
		station.Lon = e.Lon
	}

	val, ok := c.stations[e.EvaNumber]
	if !ok {
		val = &Station{EvaNumber: e.EvaNumber}
		c.stations[e.EvaNumber] = val
	}
	if e.Name != "" {
		val.Name = e.Name
	}
	if e.Lon != 0 {
		val.Lat = e.Lat
	}
	if e.Lon != 0 {
		val.Lon = e.Lon
	}
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

func (c *consumer) UpsertLineEdge(e providers.ProviderLineEdge) {
	line, ok := c.lines[e.LineID]
	if !ok {
		log.Print("Non-existant Line %d for edge upsert", e.LineID)
		return
	}
	found := false
	for _, edge := range line.Route {
		if edge.From.EvaNumber == e.EvaNumberFrom && edge.To.EvaNumber == e.EvaNumberTo {
			if e.ProviderShortestPath != nil {
				edge.ProviderShortestPath = *e.ProviderShortestPath
			}
			found = true
		}
	}
	if !found {
		log.Printf("Provider found connection that was not found by TSTP (From: %d, To: %d, LineID: %d)", e.EvaNumberFrom, e.EvaNumberTo, e.LineID)
		from, ok1 := c.stations[e.EvaNumberFrom]
		to, ok2 := c.stations[e.EvaNumberTo]
		if !ok1 || !ok2 {
			log.Panicf("Non-existant Station for stop of Line %d", e.LineID)
			return
		}
		edge := &Edge{
			Line: line,
			From: from,
			To:   to,
			Actual: StopInfo{
				Departure:      e.Planned.Departure,
				Arrival:        e.Planned.Arrival,
				DepartureTrack: e.Planned.DepartureTrack,
				ArrivalTrack:   e.Planned.ArrivalTrack,
			},
			ProviderShortestPath: *e.ProviderShortestPath,
		}
		line.Route = append(line.Route, edge)
		edge.From.Departures = append(edge.From.Departures, edge)
		edge.To.Arrivals = append(edge.To.Arrivals, edge)
	}
}

func copyProviderStopInfo(from *providers.ProviderLineStopInfo, to *StopInfo) {
	if to.Departure.IsZero() {
		to.Departure = from.Departure
		if from.Departure.IsZero() {
			to.Departure = from.Arrival
		}
		to.DepartureTrack = from.DepartureTrack
	}
	if to.Arrival.IsZero() {
		to.Arrival = from.Arrival
		if from.Arrival.IsZero() {
			to.Arrival = from.Departure
		}
		to.ArrivalTrack = from.ArrivalTrack
	}
}

func (c *consumer) initializeProviders(evaNumbers []int) {
	c.providers = []providers.Provider{&dbrest.DbRest{}}
	c.providerStations = defaultStations(evaNumbers)

	c.stations = map[int]*Station{}
	c.lines = map[int]*Line{}
}

func (c *consumer) callProviders(enrich bool) {
	for _, p := range c.providers {
		if !enrich {
			p.Fetch(c)
		} else {
			p.Enrich(c)
		}
	}
}

func defaultStations(evaNumbers []int) []providers.ProviderStation {
	var stations []providers.ProviderStation
	for _, n := range evaNumbers {
		stations = append(stations, providers.ProviderStation{EvaNumber: n})
	}
	return stations
}

func (c *consumer) rankStations(origin *Station, destination *Station) {
	var stationsSlice []*Station
	for _, s := range c.stations {
		stationsSlice = append(stationsSlice, s)
	}
	sort.Slice(stationsSlice, func(i, j int) bool {
		if stationsSlice[i] == origin || stationsSlice[j] == destination {
			return true
		}
		if stationsSlice[j] == origin || stationsSlice[i] == destination {
			return false
		}
		return geoDistStations(origin, stationsSlice[i]) < geoDistStations(origin, stationsSlice[j])
	})
	i := 0
	for _, s := range stationsSlice {
		c.stations[s.EvaNumber].Rank = i
		i++
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

func ObtainData(from int, to int, vias []int, dateTime string) (map[int]*Station, map[int]*Line) {
	c := &consumer{}

	c.parseDate(dateTime)

	var evaNumbers []int
	evaNumbers = append(evaNumbers, from)
	evaNumbers = append(evaNumbers, vias...)
	evaNumbers = append(evaNumbers, to)

	c.initializeProviders(evaNumbers)
	c.callProviders(false)
	c.generateEdges(c.stations[from], c.stations[to])
	c.rankStations(c.stations[from], c.stations[to])
	shortestPaths(c.stations, c.stations[to])
	c.callProviders(true)
	return c.stations, c.lines
}

func (c *consumer) parseDate(dateTime string) {
	layout := "2006-01-02T15:04"
	t, err := time.Parse(layout, dateTime)

	if err != nil {
		log.Panic(err)
	}
	c.dateTime = t
}
