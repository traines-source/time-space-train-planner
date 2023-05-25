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

// TODO some stations, e.g. Hamburg Hbf, yield more than 1000 results within 4 hours. Maybe filter out local transport (buses, trams etc.) in request if no vias station is specified that is nearby (reasonably reachable by local transport)?
const results = 3000

type DbRest struct {
	consumer       providers.Consumer
	client         *apiclient.Dbrest
	cachedJourneys *operations.GetJourneysOKBody
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
		r := httptransport.New(os.Getenv("API_CACHE_HOST"), os.Getenv("HAFAS_API_CACHE_PREFIX"), []string{os.Getenv("API_CACHE_SCHEME")})
		p.client = apiclient.New(r, strfmt.Default)
	}
}

func (p *DbRest) requestAndParseDeparturesArrivals() error {
	stations := p.consumer.Stations()
	for i, station := range stations {
		if i > 20 {
			log.Print("Aborting station retrieval, maximum station count exceeded.")
			break
		}
		from, to := p.consumer.RequestStationDataBetween(&station)
		duration := to.Sub(from).Minutes()
		log.Print("Requesting for ", station.EvaNumber, " at ", from, " with duration ", duration)
		if err := p.requestAndParseArrival(station, from, int64(duration)); err != nil {
			return err
		}
		if err := p.requestAndParseDeparture(station, from, int64(duration)); err != nil {
			return err
		}
	}
	return nil
}

func (p *DbRest) requestAndParseArrival(station providers.ProviderStation, when time.Time, duration int64) error {
	var params = operations.NewGetStopsIDArrivalsParams()
	params.ID = strconv.Itoa(station.EvaNumber)
	params.Duration = &duration
	r := int64(results)
	params.Results = &r
	params.When = (*strfmt.DateTime)(&when)

	res, err := p.client.Operations.GetStopsIDArrivals(params)
	if err != nil {
		return err
	}
	p.parseDepartureArrival(res.Payload, station.EvaNumber, true)
	return nil
}

func (p *DbRest) requestAndParseDeparture(station providers.ProviderStation, when time.Time, duration int64) error {
	var params = operations.NewGetStopsIDDeparturesParams()
	params.ID = strconv.Itoa(station.EvaNumber)
	params.Duration = &duration
	r := int64(results)
	params.Results = &r
	params.When = (*strfmt.DateTime)(&when)
	if station.NoLocalTransport {
		f := false
		params.Bus = &f
		params.Subway = &f
		params.Tram = &f
		params.Taxi = &f
	}

	res, err := p.client.Operations.GetStopsIDDepartures(params)
	if err != nil {
		return err
	}
	p.parseDepartureArrival(res.Payload, station.EvaNumber, false)
	return nil
}

func (p *DbRest) parseDepartureArrival(stops []*models.DepartureArrival, groupNumber int, arrival bool) {
	if len(stops) >= results {
		log.Printf("Warning: Potentially missing arrivals/departures (max. results of %d exceeded)", len(stops))
	}
	for _, stop := range stops {
		lineID, err := strconv.Atoi(*stop.Line.FahrtNr)
		if err != nil {
			log.Printf("Failed to convert Line ID %d", stop.Line.FahrtNr)
			continue
		}
		evaNumber, err := strconv.Atoi(*stop.Stop.ID)
		if err != nil {
			log.Printf("Failed to convert Eva Number %s", *stop.Stop.ID)
			continue
		}
		p.parseStation(stop, evaNumber, groupNumber)
		p.parseLine(stop, *stop.TripID, lineID)
		p.parseLineStop(stop, arrival, evaNumber, *stop.TripID)
	}
}

func (p *DbRest) parseStation(stop *models.DepartureArrival, evaNumber int, groupNumber int) {
	p.consumer.UpsertStation(providers.ProviderStation{
		EvaNumber:   evaNumber,
		GroupNumber: &groupNumber,
		Name:        *stop.Stop.Name,
		Lat:         float32(*stop.Stop.Location.Latitude),
		Lon:         float32(*stop.Stop.Location.Longitude),
	})
}

func (p *DbRest) parseLine(stop *models.DepartureArrival, tripID string, lineID int) {
	lineName := ""
	if stop.Line.Name != nil {
		lineName = *stop.Line.Name
	}
	productName := ""
	if stop.Line.Product != nil {
		productName = *stop.Line.Product
	}
	direction := ""
	if stop.Direction != nil {
		direction = *stop.Direction
	}
	p.consumer.UpsertLine(providers.ProviderLine{ID: tripID, TripName: lineID, Type: productName, Name: lineName, Direction: direction})
}

func (p *DbRest) parseLineStop(stop *models.DepartureArrival, arrival bool, evaNumber int, tripID string) {

	planned := &providers.ProviderLineStopInfo{}
	current := &providers.ProviderLineStopInfo{}

	if arrival {
		if stop.PlannedWhen != nil {
			planned.Arrival = time.Time(*stop.PlannedWhen)
		}
		if stop.PlannedPlatform != nil {
			planned.ArrivalTrack = *stop.PlannedPlatform
		}
		if stop.When != nil && stop.Delay != nil {
			current.Arrival = time.Time(*stop.When)
		}
		if stop.Platform != nil {
			current.ArrivalTrack = *stop.Platform
		}
	} else {
		if stop.PlannedWhen != nil {
			planned.Departure = time.Time(*stop.PlannedWhen)
		}
		if stop.PlannedPlatform != nil {
			planned.DepartureTrack = *stop.PlannedPlatform
		}
		if stop.When != nil && stop.Delay != nil {
			current.Departure = time.Time(*stop.When)
		}
		if stop.Platform != nil {
			current.DepartureTrack = *stop.Platform
		}
	}
	pls := providers.ProviderLineStop{EvaNumber: evaNumber, LineID: tripID, Planned: planned, Current: current, Cancelled: stop.Cancelled}
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
	from := strconv.Itoa(stations[0].EvaNumber)
	to := strconv.Itoa(stations[len(stations)-1].EvaNumber)
	params.From = &from
	params.To = &to
	params.Departure = (*strfmt.DateTime)(&departure)
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
			evaNumberFrom, err1 := strconv.Atoi(*leg.Origin.ID)
			evaNumberTo, err2 := strconv.Atoi(*leg.Destination.ID)
			if err1 == nil && err2 == nil {
				from := providers.ProviderStation{
					EvaNumber: evaNumberFrom,
					Name:      *leg.Origin.Name,
					Lat:       float32(*leg.Origin.Location.Latitude),
					Lon:       float32(*leg.Origin.Location.Longitude),
				}
				to := providers.ProviderStation{
					EvaNumber: evaNumberTo,
					Name:      *leg.Destination.Name,
					Lat:       float32(*leg.Destination.Location.Latitude),
					Lon:       float32(*leg.Destination.Location.Longitude),
				}
				p.fallbackStations(from, to)
				p.consumer.UpsertStation(from)
				p.consumer.UpsertStation(to)
				if leg.Arrival != nil && end.Before(time.Time(*leg.Arrival)) {
					end = time.Time(*leg.Arrival)
				}
			} else {
				log.Print("Error while trying to read stations from journeys")
			}
		}
	}
	start, _ := p.consumer.RequestStationDataBetween(&p.consumer.Stations()[0])
	log.Print("expdur", start, end)
	p.consumer.SetExpectedTravelDuration(end.Sub(start))
}


func (p *DbRest) fallbackStations(from providers.ProviderStation, to providers.ProviderStation) {
	stations := p.consumer.Stations()
	p.consumer.UpsertStation(providers.ProviderStation{
		EvaNumber:   stations[0].EvaNumber,
		GroupNumber: &stations[0].EvaNumber,
		Lat:         from.Lat,
		Lon:         from.Lon,
	})
	p.consumer.UpsertStation(providers.ProviderStation{
		EvaNumber:   stations[len(stations)-1].EvaNumber,
		GroupNumber: &stations[len(stations)-1].EvaNumber,
		Lat:         to.Lat,
		Lon:         to.Lon,
	})
}

func (p *DbRest) parseEdgesFromJourneys() {
	for _, journey := range p.cachedJourneys.Journeys {
		for _, leg := range journey.Legs {
			evaNumberFrom, err1 := strconv.Atoi(*leg.Origin.ID)
			evaNumberTo, err2 := strconv.Atoi(*leg.Destination.ID)
			if err1 != nil || err2 != nil || leg.Line == nil {
				log.Print("Error while trying to read edges from journeys ", err1, err2, leg.Line)
				continue
			}
			hafas := true
			planned := &providers.ProviderLineStopInfo{}
			if leg.Departure != nil {
				planned.Departure = time.Time(*leg.Departure)
			}
			if leg.Arrival != nil {
				planned.Arrival = time.Time(*leg.Arrival)
			}
			p.consumer.UpsertLineEdge(providers.ProviderLineEdge{
				EvaNumberFrom:        evaNumberFrom,
				EvaNumberTo:          evaNumberTo,
				LineID:               *leg.TripID,
				ProviderShortestPath: &hafas,
				Planned:              planned,
			})
		}
	}
}
