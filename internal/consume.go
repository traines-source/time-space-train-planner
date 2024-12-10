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
	stations               map[string]*Station
	lines                  map[string]*Line
	dateTime               time.Time
	regionalOnly           bool
	expectedTravelDuration time.Duration
}

var loc, _ = time.LoadLocation("Europe/Berlin")

func (c *consumer) RequestStationDataBetween(station *providers.ProviderStation) (from time.Time, to time.Time) {
	defaultDuration, _ := time.ParseDuration("1h")
	longDuration, _ := time.ParseDuration("8h")
	maxDuration, _ := time.ParseDuration("10h")

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

func (c *consumer) StationByID(id string) (providers.ProviderStation, error) {
	for _, v := range c.providerStations {
		if id == v.ID {
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
		if s.ID == e.ID {
			station = &s
		}
	}
	if station == nil {
		station = &providers.ProviderStation{ID: e.ID}
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
	if e.GroupID != nil {
		station.GroupID = e.GroupID
	}

	val, ok := c.stations[e.ID]
	if !ok {
		val = &Station{ID: e.ID}
		c.stations[e.ID] = val
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
	if e.GroupID != nil {
		val.GroupID = e.GroupID
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
	station, ok := c.stations[e.ID]
	if !ok {
		log.Panicf("Non-existant Station %s for stop of Line %s", e.ID, e.LineID)
		return
	}
	line, ok := c.lines[e.LineID]
	if !ok {
		log.Panicf("Non-existant Line %s for Station  %s", e.LineID, e.ID)
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
		log.Printf("Provider found Line that was not found by TSTP (From: %s, To: %s, LineID: %s, Dep: %s)", e.IDFrom, e.IDTo, e.LineID, e.Planned.Departure)
		return
	}
	foundStart := false
	foundEnd := false
	for _, edge := range line.Route {
		// TODO handle multi-line trains (ICE / RE, IC / NJ etc, e.g. IC 60400/NJ 40470)
		if edge.From.ID == e.IDFrom || c.sameGroupNumber(e.IDFrom, edge.From.GroupID) || foundStart && !foundEnd {
			if e.ProviderShortestPath != nil {
				edge.ProviderShortestPath = *e.ProviderShortestPath
			}
			foundStart = true
			if edge.To.ID == e.IDTo || c.sameGroupNumber(e.IDTo, edge.To.GroupID) {
				foundEnd = true
			}
		}
	}
	if !foundEnd {
		log.Printf("Provider found connection that was not found by TSTP (From: %s, To: %s, LineID: %s, Name: %s, Dep: %s, foundStart: %t)", e.IDFrom, e.IDTo, e.LineID, line.Name, e.Planned.Departure, foundStart)
	}
}

func (c *consumer) sameGroupNumber(id string, groupID *string) bool {
	if val, ok := c.stations[id]; ok && groupID != nil && val.GroupID != nil && *val.GroupID == *groupID {
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

func (c *consumer) initializeProviders(stationIDs []string) {
	c.stations = map[string]*Station{}
	c.lines = map[string]*Line{}

	c.providers = []providers.Provider{&dbrest.DbRest{}}
	c.providerStations = c.defaultStations(stationIDs)
}

func (c *consumer) callProviders(call func(providers.Provider, *consumer) error) *ErrorCode {
	for _, p := range c.providers {
		if err := call(p, c); err != nil {
			return HandleError(err)
		}
	}
	log.Println("Provider requests completed.")
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

func (c *consumer) defaultStations(stationIDs []string) []providers.ProviderStation {
	var stations []providers.ProviderStation
	for _, n := range stationIDs {
		s := providers.ProviderStation{ID: n}
		stations = append(stations, s)
		c.UpsertStation(s)
	}
	return stations
}

func indexOf(slice []string, value string) int {
	for i, e := range slice {
		if e == value {
			return i
		}
	}
	return -1
}

func (c *consumer) rankStations(origin *Station, destination *Station) {
	//force := []int{8070003, 8070004, 8000105, 8098105, 8006404, 8000615}
	force := []string{}
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
		if stationsSlice[i].GroupID != nil && stationsSlice[j].GroupID != nil && *stationsSlice[i].GroupID == *stationsSlice[j].GroupID {
			return false
		}
		forceI := indexOf(force, stationsSlice[i].ID)
		forceJ := indexOf(force, stationsSlice[j].ID)
		if forceI != -1 && forceJ != -1 {
			return forceI < forceJ
		}
		stationI := stationsSlice[i]
		if stationI.GroupID != nil {
			if val, ok := c.stations[*stationI.GroupID]; ok {
				stationI = val
			}
		}
		stationJ := stationsSlice[j]
		if stationJ.GroupID != nil {
			if val, ok := c.stations[*stationJ.GroupID]; ok {
				stationJ = val
			}
		}
		a := geoDistStations(origin, stationI) - geoDistStations(destination, stationI)
		b := geoDistStations(origin, stationJ) - geoDistStations(destination, stationJ)
		return a < b
	})
	i := 0
	for _, s := range stationsSlice {
		c.stations[s.ID].Rank = i
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

func prepare(from string, to string, vias []string, dateTime string, regionly bool) *consumer {
	c := &consumer{}

	c.parseDate(dateTime)
	c.dateTime = c.dateTime.Add(-time.Minute * 5)
	c.regionalOnly = regionly

	var stationIDs []string
	stationIDs = append(stationIDs, from)
	stationIDs = append(stationIDs, vias...)
	stationIDs = append(stationIDs, to)

	log.Print(stationIDs)
	c.initializeProviders(stationIDs)
	return c
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

func (c *consumer) apiFlow(system string, from string, to string, vias []string, regionly bool) *ErrorCode {
	if err := c.callProviders(callDeparturesArrivals); err != nil {
		return err
	}
	if err := c.generateEdges(c.stations[from], c.stations[to]); err != nil {
		return err
	}
	if err := c.callProviders(callEnrich); err != nil {
		return err
	}
	StostEnrich(system, c.lines, c.stations, from, to, c.dateTime, time.Now(), regionly)
	return nil
}

func ObtainVias(from string, to string, vias []string, dateTime string, regionly bool) (map[string]*Station, *ErrorCode) {
	c := prepare(from, to, vias, dateTime, regionly)
	if err := c.callProviders(callVias); err != nil {
		return nil, err
	}
	return c.stations, nil
}

func ObtainData(system string, from string, to string, vias []string, dateTime string, regionly bool) (map[string]*Station, map[string]*Line, *ErrorCode) {
	c := prepare(from, to, vias, dateTime, regionly)
	if system == "" {
		if err := c.apiFlow(system, from, to, vias, regionly); err != nil {
			return nil, nil, err
		}
	} else {
		StostProduce(system, c.lines, c.stations, from, to, c.dateTime, time.Now())
	}
	c.rankStations(c.stations[from], c.stations[to])
	return c.stations, c.lines, nil
}
