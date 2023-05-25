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
	providers              []providers.Provider
	providerStations       []providers.ProviderStation
	stations               map[int]*Station
	lines                  map[string]*Line
	dateTime               time.Time
	regionalOnly           bool
	expectedTravelDuration time.Duration
}

var loc, _ = time.LoadLocation("Europe/Berlin")

func (c *consumer) RequestStationDataBetween(station *providers.ProviderStation) (from time.Time, to time.Time) {
	defaultDuration, _ := time.ParseDuration("2h")
	longDuration, _ := time.ParseDuration("8h")
	maxDuration, _ := time.ParseDuration("14h")

	var travelDuration time.Duration
	if c.expectedTravelDuration > maxDuration {
		travelDuration = maxDuration + defaultDuration + defaultDuration
	} else {
		travelDuration = c.expectedTravelDuration.Round(time.Hour) + defaultDuration
		if travelDuration > longDuration {
			travelDuration += defaultDuration
		}
	}
	//t := time.Now()
	//from = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.Local)
	//from = time.Date(t.Year(), t.Month(), 9, 19, 0, 0, 0, time.Local)
	from = c.dateTime
	return from, from.Add(travelDuration)
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

func (c *consumer) RegionalOnly() bool {
	return c.regionalOnly
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
	if e.GroupNumber != nil {
		station.GroupNumber = e.GroupNumber
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
	if e.GroupNumber != nil {
		val.GroupNumber = e.GroupNumber
	}
}

func (c *consumer) UpsertLine(e providers.ProviderLine) {
	val, ok := c.lines[e.ID]
	if !ok {
		val = &Line{ID: e.ID, Stops: []*LineStop{}}
		c.lines[e.ID] = val
	}
	val.Name = e.Name
	val.Type = e.Type
	val.Message = e.Message
	if e.Direction != "" {
		val.Direction = e.Direction
	}
}

func existingStopHasDifferentPlanned(e providers.ProviderLineStop, stop *LineStop) bool {
	return e.Planned != nil &&
		(!e.Planned.Arrival.IsZero() && !stop.Planned.Arrival.IsZero() && e.Planned.Arrival != stop.Planned.Arrival ||
			!e.Planned.Departure.IsZero() && !stop.Planned.Departure.IsZero() && e.Planned.Departure != stop.Planned.Departure)
}

func (c *consumer) UpsertLineStop(e providers.ProviderLineStop) {
	station, ok := c.stations[e.EvaNumber]
	if !ok {
		log.Panicf("Non-existant Station %d for stop of Line %s", e.EvaNumber, e.LineID)
		return
	}
	line, ok := c.lines[e.LineID]
	if !ok {
		log.Panicf("Non-existant Line %s for Station  %d", e.LineID, e.EvaNumber)
		return
	}
	var val *LineStop
	for _, stop := range line.Stops {
		if stop.Station == station && !existingStopHasDifferentPlanned(e, stop) {
			val = stop
			break
		}
	}
	if val == nil {
		val = &LineStop{Station: station}
		line.Stops = append(line.Stops, val)
	}
	if e.Planned != nil {
		copyProviderStopInfo(e.Planned, &val.Planned)
	}
	if e.Current != nil {
		copyProviderStopInfo(e.Current, &val.Current)
	}
	val.Message = e.Message
	val.Cancelled = e.Cancelled
}

func (c *consumer) UpsertLineEdge(e providers.ProviderLineEdge) {
	line, ok := c.lines[e.LineID]
	if !ok {
		log.Printf("Provider found Line that was not found by TSTP (From: %d, To: %d, LineID: %s, Dep: %s)", e.EvaNumberFrom, e.EvaNumberTo, e.LineID, e.Planned.Departure)
		return
	}
	foundStart := false
	foundEnd := false
	for _, edge := range line.Route {
		// TODO handle multi-line trains (ICE / RE, IC / NJ etc, e.g. IC 60400/NJ 40470)
		if edge.From.EvaNumber == e.EvaNumberFrom || c.sameGroupNumber(e.EvaNumberFrom, edge.From.GroupNumber) || foundStart && !foundEnd {
			if e.ProviderShortestPath != nil {
				edge.ProviderShortestPath = *e.ProviderShortestPath
			}
			foundStart = true
			if edge.To.EvaNumber == e.EvaNumberTo || c.sameGroupNumber(e.EvaNumberTo, edge.To.GroupNumber) {
				foundEnd = true
			}
		}
	}
	if !foundEnd {
		log.Printf("Provider found connection that was not found by TSTP (From: %d, To: %d, LineID: %s, Name: %s, Dep: %s, foundStart: %t)", e.EvaNumberFrom, e.EvaNumberTo, e.LineID, line.Name, e.Planned.Departure, foundStart)
	}
}

func (c *consumer) sameGroupNumber(evaNumber int, groupNumber *int) bool {
	if val, ok := c.stations[evaNumber]; ok && groupNumber != nil && val.GroupNumber != nil && *val.GroupNumber == *groupNumber {
		return true
	}
	return false
}

func (c *consumer) SetExpectedTravelDuration(duration time.Duration) {
	c.expectedTravelDuration = duration
}

func copyProviderStopInfo(from *providers.ProviderLineStopInfo, to *StopInfo) {
	if to.Departure.IsZero() || to.Departure == to.Arrival {
		to.Departure = from.Departure
		if from.Departure.IsZero() && !from.Arrival.IsZero() {
			to.Departure = from.Arrival
		}
		to.DepartureTrack = from.DepartureTrack
	}
	if to.Arrival.IsZero() || to.Departure == to.Arrival {
		to.Arrival = from.Arrival
		if from.Arrival.IsZero() && !from.Departure.IsZero() {
			to.Arrival = from.Departure
		}
		to.ArrivalTrack = from.ArrivalTrack
	}
}

func (c *consumer) initializeProviders(evaNumbers []int) {
	c.stations = map[int]*Station{}
	c.lines = map[string]*Line{}

	c.providers = []providers.Provider{&dbrest.DbRest{}}
	c.providerStations = c.defaultStations(evaNumbers)
}

func (c *consumer) callProviders(call func(providers.Provider, *consumer) error) *ErrorCode {
	for _, p := range c.providers {
		if err := call(p, c); err != nil {
			return HandleError(err)
		}
	}
	return nil
}

func callVias(p providers.Provider, c *consumer) error {
	return p.Vias(c)
}

func callDeparturesArrivals(p providers.Provider, c *consumer) error {
	return p.DeparturesArrivals(c)
}

func callEnrich(p providers.Provider, c *consumer) error {
	return p.Enrich(c)
}

func (c *consumer) defaultStations(evaNumbers []int) []providers.ProviderStation {
	var stations []providers.ProviderStation
	for _, n := range evaNumbers {
		s := providers.ProviderStation{EvaNumber: n}
		stations = append(stations, s)
		c.UpsertStation(s)
	}
	return stations
}

func indexOf(slice []int, value int) int {
	for i, e := range slice {
		if e == value {
			return i
		}
	}
	return -1
}

func (c *consumer) rankStations(origin *Station, destination *Station) {
	//force := []int{8070003, 8070004, 8000105, 8098105, 8006404, 8000615}
	force := []int{}
	var stationsSlice []*Station
	for _, s := range c.stations {
		stationsSlice = append(stationsSlice, s)
	}
	sort.SliceStable(stationsSlice, func(i, j int) bool {
		return stationsSlice[i].Name < stationsSlice[j].Name
	})
	sort.SliceStable(stationsSlice, func(i, j int) bool {
		if stationsSlice[i] == origin || stationsSlice[j] == destination {
			return true
		}
		if stationsSlice[j] == origin || stationsSlice[i] == destination {
			return false
		}
		if stationsSlice[i].GroupNumber != nil && stationsSlice[j].GroupNumber != nil && *stationsSlice[i].GroupNumber == *stationsSlice[j].GroupNumber {
			return false
		}
		forceI := indexOf(force, stationsSlice[i].EvaNumber)
		forceJ := indexOf(force, stationsSlice[j].EvaNumber)
		if forceI != -1 && forceJ != -1 {
			return forceI < forceJ
		}
		stationI := stationsSlice[i]
		if stationI.GroupNumber != nil {
			if val, ok := c.stations[*stationI.GroupNumber]; ok {
				stationI = val
			}
		}
		stationJ := stationsSlice[j]
		if stationJ.GroupNumber != nil {
			if val, ok := c.stations[*stationJ.GroupNumber]; ok {
				stationJ = val
			}
		}
		a := geoDistStations(origin, stationI) - geoDistStations(destination, stationI)
		b := geoDistStations(origin, stationJ) - geoDistStations(destination, stationJ)
		return a < b
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
	if thisFrom.ArrivalTrack != "" {
		to.ArrivalTrack = thisFrom.ArrivalTrack
	}
	if !lastFrom.Departure.IsZero() {
		to.Departure = lastFrom.Departure
	}
	if !thisFrom.Arrival.IsZero() {
		to.Arrival = thisFrom.Arrival
	}
}

func prepare(from int, to int, vias []int, dateTime string, regionly bool) *consumer {
	c := &consumer{}

	c.parseDate(dateTime)
	c.regionalOnly = regionly

	var evaNumbers []int
	evaNumbers = append(evaNumbers, from)
	evaNumbers = append(evaNumbers, vias...)
	evaNumbers = append(evaNumbers, to)

	log.Print(evaNumbers)
	c.initializeProviders(evaNumbers)
	return c
}

func ObtainVias(from int, to int, vias []int, dateTime string, regionly bool) (map[int]*Station, *ErrorCode) {
	c := prepare(from, to, vias, dateTime, regionly)
	if err := c.callProviders(callVias); err != nil {
		return nil, err
	}
	return c.stations, nil
}

func ObtainData(from int, to int, vias []int, dateTime string, regionly bool) (map[int]*Station, map[string]*Line, *ErrorCode) {
	c := prepare(from, to, vias, dateTime, regionly)
	if err := c.callProviders(callDeparturesArrivals); err != nil {
		return nil, nil, err
	}
	if err := c.generateEdges(c.stations[from], c.stations[to]); err != nil {
		return nil, nil, err
	}
	shortestPaths(c.stations, c.stations[from], c.stations[to], regionly)
	if err := c.callProviders(callEnrich); err != nil {
		return nil, nil, err
	}
	c.rankStations(c.stations[from], c.stations[to])
	return c.stations, c.lines, nil
}

func (c *consumer) parseDate(dateTime string) {
	layout := "2006-01-02T15:04"
	t, err := time.ParseInLocation(layout, dateTime, loc)

	if err != nil {
		t := time.Now()
		c.dateTime = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.Local)
	} else {
		c.dateTime = t
	}
}
