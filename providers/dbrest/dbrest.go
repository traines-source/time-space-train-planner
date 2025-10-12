package dbrest

import (
	"log"
	"os"
	"strconv"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"traines.eu/time-space-train-planner/providers"
	apiclient "traines.eu/time-space-train-planner/providers/dbrest/client"
	"traines.eu/time-space-train-planner/providers/dbrest/client/operations"
	"traines.eu/time-space-train-planner/providers/dbrest/models"
)

const defaultResults = 3000

type DbRest struct {
	consumer       providers.Consumer
	client         *apiclient.Dbrest
	cachedJourneys *operations.GetJourneysOKBody
	backend        string
}

// TODO some stations, e.g. Hamburg Hbf, yield more than 1000 defaultResults within 4 hours. Maybe filter out local transport (buses, trams etc.) in request if no vias station is specified that is nearby (reasonably reachable by local transport)?
func (p *DbRest) results(seen int) int {
	if p.backend == "transitous" && seen >= defaultResults {
		return defaultResults * 3
	}
	return defaultResults
}

func (p *DbRest) SetBackend(backend string) {
	p.backend = backend
}

func (p *DbRest) Vias(c providers.Consumer) error {
	p.prepareClient(c)
	if err := p.requestJourneys(); err != nil {
		return err
	}
	p.parseStationsFromJourneys()
	return nil
}

func (p *DbRest) DeparturesArrivals(c providers.Consumer) error {
	p.prepareClient(c)
	if err := p.requestJourneys(); err != nil {
		return err
	}
	p.parseStationsFromJourneys()
	if err := p.requestAndParseDeparturesArrivals(); err != nil {
		return err
	}
	return nil
}

func (p *DbRest) Enrich(c providers.Consumer) error {
	p.prepareClient(c)
	if err := p.requestJourneys(); err != nil {
		return err
	}
	p.parseEdgesFromJourneys()
	return nil
}

func (p *DbRest) prepareClient(c providers.Consumer) {
	p.consumer = c
	if p.client == nil {
		prefix := os.Getenv("HAFAS_API_CACHE_PREFIX")
		if p.backend == "transitous" {
			prefix = os.Getenv("TRANSITOUS_API_CACHE_PREFIX")
		}
		log.Println(os.Getenv("API_CACHE_HOST"))
		r := httptransport.New(os.Getenv("API_CACHE_HOST"), prefix, []string{os.Getenv("API_CACHE_SCHEME")})
		p.client = apiclient.New(r, strfmt.Default)
	}
}

func (p *DbRest) requestAndParseDeparturesArrivals() error {
	stations := p.consumer.Stations()
	seen := map[string]int{}
	i := 0
	for _, station := range stations {
		if seen[station.ID] != 0 {
			log.Print("Skipping already seen ", station.ID, " with ", seen[station.ID])
			continue
		}
		if i > 30 {
			log.Print("Aborting station retrieval, maximum station count exceeded.")
			break
		}
		i++
		from, to := p.consumer.RequestStationDataBetween(&station)
		duration := to.Sub(from).Minutes()
		log.Print("Requesting for ", station.ID, " at ", from, " with duration ", duration, " numResults ", p.results(seen[station.ID]))
		if err := p.requestAndParseArrival(station, from, int64(duration), seen); err != nil {
			return err
		}
		if err := p.requestAndParseDeparture(station, from, int64(duration), seen); err != nil {
			return err
		}
	}
	return nil
}

func (p *DbRest) requestAndParseArrival(station providers.ProviderStation, when time.Time, duration int64, seen map[string]int) error {
	var params = operations.NewGetStopsIDArrivalsParams()
	params.ID = station.ID
	params.Duration = &duration
	r := int64(p.results(seen[station.ID]))
	params.Results = &r
	params.When = (*strfmt.DateTime)(&when)

	for {
		res, err := p.client.Operations.GetStopsIDArrivals(params)
		if err != nil {
			return err
		}
		r := int64(p.results(len(res.Payload.Arrivals)))
		if len(res.Payload.Arrivals) < defaultResults || *params.Results >= r {
			p.parseDepartureArrival(res.Payload.Arrivals, station.ID, true, seen)
			break
		}
		params.Results = &r
		log.Print("Requesting again for ", station.ID, " with duration ", duration, " numResults ", r, " ", len(res.Payload.Arrivals))
	}
	return nil
}

func (p *DbRest) requestAndParseDeparture(station providers.ProviderStation, when time.Time, duration int64, seen map[string]int) error {
	var params = operations.NewGetStopsIDDeparturesParams()
	params.ID = station.ID
	params.Duration = &duration
	r := int64(p.results(seen[station.ID]))
	params.Results = &r
	params.When = (*strfmt.DateTime)(&when)
	if station.NoLocalTransport {
		f := false
		params.Bus = &f
		params.Subway = &f
		params.Tram = &f
		params.Taxi = &f
	}

	for {
		res, err := p.client.Operations.GetStopsIDDepartures(params)
		if err != nil {
			return err
		}
		r := int64(p.results(len(res.Payload.Departures)))
		if len(res.Payload.Departures) < defaultResults || *params.Results >= r {
			p.parseDepartureArrival(res.Payload.Departures, station.ID, false, seen)
			break
		}
		params.Results = &r
		log.Print("Requesting again for ", station.ID, " with duration ", duration, " numResults ", r, " ", len(res.Payload.Departures))
	}
	return nil
}

func (p *DbRest) parseDepartureArrival(stops []*models.Alternative, groupID string, arrival bool, seen map[string]int) {
	if len(stops) >= defaultResults {
		log.Printf("Warning: Potentially missing arrivals/departures (max. results of %d exceeded for %s)", len(stops), groupID)
	}
	for _, stop := range stops {
		lineID, err := strconv.Atoi(stop.Line.FahrtNr)
		if err != nil {
			//log.Printf("Failed to convert Line ID %s", stop.Line.FahrtNr)
			lineID = 0
		}
		tripID := getNormalizedTripID(stop.TripID, stop.Line.ID, stop.Line.FahrtNr, stop.Line.ProductName)
		p.parseStation(stop.Stop, stop.Stop.ID, groupID)
		seen[stop.Stop.ID] = len(stops)
		p.parseLine(stop, tripID, lineID)
		p.parseLineStop(stop, arrival, stop.Stop.ID, tripID)
	}
}
func getNormalizedTripID(tripID string, lineID string, fahrtNr string, productName string) string {
	/*if len(lineID) >= 4 && len(fahrtNr) >= 3 && productName != "Bus" {
		parts := strings.Split(tripID, "#")
		if matched, _ := regexp.MatchString("[0-9]", lineID); matched && len(parts) == 42 {
			date := parts[12]
			id := lineID + "###" + fahrtNr + "###" + date
			return id
		}
	}*/
	return tripID
}

func (p *DbRest) parseStation(stop *models.Stop, stationID string, groupID string) {
	s := providers.ProviderStation{
		ID:      stationID,
		GroupID: &groupID,
		Name:    stop.Name,
	}
	if stop.Location != nil {
		s.Lat = stop.Location.Latitude
		s.Lon = stop.Location.Longitude
	}
	p.consumer.UpsertStation(s)
}

func (p *DbRest) parseLine(stop *models.Alternative, tripID string, lineID int) {
	lineName := ""
	if stop.Line.Name != "" {
		lineName = stop.Line.Name
	}
	productName := ""
	if stop.Line.Product != "" {
		productName = stop.Line.Product
	}
	direction := ""
	if stop.Direction != "" {
		direction = stop.Direction
	}
	p.consumer.UpsertLine(providers.ProviderLine{ID: tripID, TripName: lineID, Type: productName, Name: lineName, Direction: direction})
}

func (p *DbRest) parseLineStop(stop *models.Alternative, arrival bool, stationID string, tripID string) {

	planned := &providers.ProviderLineStopInfo{}
	current := &providers.ProviderLineStopInfo{}

	if arrival {
		if !stop.PlannedWhen.IsZero() {
			planned.Arrival = time.Time(stop.PlannedWhen)
		}
		if stop.PlannedPlatform != "" {
			planned.ArrivalTrack = stop.PlannedPlatform
		}
		if !stop.When.IsZero() && stop.Delay != nil {
			current.Arrival = time.Time(stop.When)
		}
		if stop.Platform != "" {
			current.ArrivalTrack = stop.Platform
		}
	} else {
		if !stop.PlannedWhen.IsZero() {
			planned.Departure = time.Time(stop.PlannedWhen)
		}
		if stop.PlannedPlatform != "" {
			planned.DepartureTrack = stop.PlannedPlatform
		}
		if !stop.When.IsZero() && stop.Delay != nil {
			current.Departure = time.Time(stop.When)
		}
		if stop.Platform != "" {
			current.DepartureTrack = stop.Platform
		}
	}
	pls := providers.ProviderLineStop{ID: stationID, LineID: tripID, Planned: planned, Current: current, Cancelled: stop.Cancelled}
	if len(stop.Remarks) > 0 {
		for _, remark := range stop.Remarks {
			if pls.Message != "" {
				pls.Message += ", " + remark.Text
			} else {
				pls.Message = remark.Text
			}
		}
	}
	p.consumer.UpsertLineStop(pls)
}

func (p *DbRest) requestJourneys() error {
	if p.cachedJourneys == nil {
		if err := p.requestJourneysApi(); err != nil {
			return err
		}
	}
	return nil
}

func (p *DbRest) requestJourneysApi() error {
	stations := p.consumer.Stations()
	departure, _ := p.consumer.RequestStationDataBetween(&stations[0])
	var params = operations.NewGetJourneysParams()
	from := stations[0].ID
	to := stations[len(stations)-1].ID
	params.From = &from
	params.To = &to
	params.Departure = (*strfmt.DateTime)(&departure)
	var resultNum int64 = 10
	params.Results = &resultNum
	if p.consumer.RegionalOnly() {
		falsey := false
		params.National = &falsey
		params.NationalExpress = &falsey
	}
	res, err := p.client.Operations.GetJourneys(params)
	if err != nil {
		return err
	}
	p.cachedJourneys = res.Payload
	return nil
}

func (p *DbRest) parseStationsFromJourneys() {
	var end time.Time
	for _, journey := range p.cachedJourneys.Journeys {
		for _, leg := range journey.Legs {
			p.parseStation(leg.Origin, leg.Origin.ID, leg.Origin.ID)
			p.parseStation(leg.Destination, leg.Destination.ID, leg.Destination.ID)
			if !leg.Arrival.IsZero() && end.Before(time.Time(leg.Arrival)) {
				end = time.Time(leg.Arrival)
			}
		}
		p.fallbackStations(journey)
	}
	start, _ := p.consumer.RequestStationDataBetween(&p.consumer.Stations()[0])
	p.consumer.SetExpectedTravelDuration(end.Sub(start))
}

func (p *DbRest) fallbackStations(j *models.Journey) {
	stations := p.consumer.Stations()
	if len(j.Legs) == 0 {
		return
	}
	p.parseStation(j.Legs[0].Origin, stations[0].ID, stations[0].ID)
	p.parseStation(j.Legs[len(j.Legs)-1].Destination, stations[len(stations)-1].ID, stations[len(stations)-1].ID)
}

func (p *DbRest) parseEdgesFromJourneys() {
	log.Println("Num Journeys: ", len(p.cachedJourneys.Journeys))
	for _, journey := range p.cachedJourneys.Journeys {
		for _, leg := range journey.Legs {
			if leg.Line == nil {
				if !leg.Walking {
					log.Print("Error while trying to read edges from journeys ", leg.Line)
				}
				continue
			}
			hafas := true
			planned := &providers.ProviderLineStopInfo{}
			if !leg.Departure.IsZero() {
				planned.Departure = time.Time(leg.Departure)
			}
			if !leg.Arrival.IsZero() {
				planned.Arrival = time.Time(leg.Arrival)
			}
			p.consumer.UpsertLineEdge(providers.ProviderLineEdge{
				IDFrom:               leg.Origin.ID,
				IDTo:                 leg.Destination.ID,
				LineID:               getNormalizedTripID(leg.TripID, leg.Line.ID, leg.Line.FahrtNr, leg.Line.ProductName),
				ProviderShortestPath: &hafas,
				Planned:              planned,
			})
		}
	}
}
