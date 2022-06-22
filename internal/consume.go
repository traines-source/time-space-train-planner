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
	expectedTravelDuration time.Duration
}

var loc, _ = time.LoadLocation("Europe/Berlin")

func (c *consumer) RequestStationDataBetween(station *providers.ProviderStation) (from time.Time, to time.Time) {
	// TODO increase depending on journey time according to HAFAS, otherwise longer journeys are impossible to plan
	defaultDuration, _ := time.ParseDuration("2h")
	maxDuration, _ := time.ParseDuration("10h")

	var travelDuration time.Duration
	if c.expectedTravelDuration < defaultDuration {
		travelDuration = defaultDuration
	} else if c.expectedTravelDuration > maxDuration {
		travelDuration = maxDuration
	} else {
		travelDuration = c.expectedTravelDuration.Round(time.Hour)
	}
	log.Print("Requesting for ", station.EvaNumber, " at ", c.dateTime, " with duration +2 ", travelDuration)
	//t := time.Now()
	//from = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.Local)
	//from = time.Date(t.Year(), t.Month(), 9, 19, 0, 0, 0, time.Local)
	from = c.dateTime
	return from, from.Add(travelDuration).Add(defaultDuration)
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
}

func (c *consumer) UpsertLineEdge(e providers.ProviderLineEdge) {
	line, ok := c.lines[e.LineID]
	if !ok {
		log.Printf("Non-existant Line %s for edge upsert", e.LineID)
		return
	}
	foundStart := false
	foundEnd := false
	for _, edge := range line.Route {
		// TODO handle multi-line trains (ICE / RE, IC / NJ etc, e.g. IC 60400/NJ 40470)
		if edge.From.EvaNumber == e.EvaNumberFrom || foundStart && !foundEnd {
			if e.ProviderShortestPath != nil {
				edge.ProviderShortestPath = *e.ProviderShortestPath
			}
			foundStart = true
			if edge.To.EvaNumber == e.EvaNumberTo {
				foundEnd = true
			}
		}
	}
	if !foundEnd {
		log.Printf("Provider found connection that was not found by TSTP (From: %d, To: %d, LineID: %s)", e.EvaNumberFrom, e.EvaNumberTo, e.LineID)
	}
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
	c.providers = []providers.Provider{&dbrest.DbRest{}}
	c.providerStations = defaultStations(evaNumbers)

	c.stations = map[int]*Station{}
	c.lines = map[string]*Line{}
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
			stationI = c.stations[*stationI.GroupNumber]
		}
		stationJ := stationsSlice[j]
		if stationJ.GroupNumber != nil {
			stationJ = c.stations[*stationJ.GroupNumber]
		}
		return geoDistStations(origin, stationI)-geoDistStations(destination, stationI) < geoDistStations(origin, stationJ)-geoDistStations(destination, stationJ)
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

func ObtainData(from int, to int, vias []int, dateTime string) (map[int]*Station, map[string]*Line) {
	c := &consumer{}

	c.parseDate(dateTime)

	var evaNumbers []int
	evaNumbers = append(evaNumbers, from)
	evaNumbers = append(evaNumbers, vias...)
	evaNumbers = append(evaNumbers, to)

	log.Print(evaNumbers)
	c.initializeProviders(evaNumbers)
	c.callProviders(false)
	c.generateEdges(c.stations[from], c.stations[to])
	shortestPaths(c.stations, c.stations[from], c.stations[to])
	c.callProviders(true)
	c.rankStations(c.stations[from], c.stations[to])
	return c.stations, c.lines
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
