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

type DbRest struct {
	consumer providers.Consumer
	client   *apiclient.Dbrest
}

func (p *DbRest) Fetch(c providers.Consumer) {
	p.consumer = c
	p.prepareClient()
	//p.requestStations()
	p.requestDeparturesAndArrivals()
}

func (p *DbRest) requestStations() {
	stations := p.consumer.Stations()
	for _, station := range stations {
		p.requestStation(station)
	}
}

func (p *DbRest) requestDeparturesAndArrivals() {
	stations := p.consumer.Stations()
	for _, station := range stations {
		from, to := p.consumer.RequestStationDataBetween(&station)
		duration := to.Sub(from).Minutes()
		p.requestArrival(station, from, int64(duration))
		p.requestDeparture(station, from, int64(duration))
	}
}

func (p *DbRest) requestStation(station providers.ProviderStation) {
	var params = operations.NewGetStationsIDParams()
	params.ID = strconv.Itoa(station.EvaNumber)
	res, err := p.client.Operations.GetStationsID(params)
	if err != nil {
		log.Panic(err)
		return
	}
	p.consumer.UpsertStation(providers.ProviderStation{
		EvaNumber: station.EvaNumber,
		Lat:       float32(res.Payload.Location.Latitude),
		Lon:       float32(res.Payload.Location.Longitude),
	})
}

func (p *DbRest) requestArrival(station providers.ProviderStation, when time.Time, duration int64) {
	var params = operations.NewGetStopsIDArrivalsParams()
	params.ID = strconv.Itoa(station.EvaNumber)
	params.Duration = &duration
	params.When = (*strfmt.DateTime)(&when)

	res, err := p.client.Operations.GetStopsIDArrivals(params)
	if err != nil {
		log.Panic(err)
		return
	}
	p.consumer.UpsertStation(providers.ProviderStation{
		EvaNumber: station.EvaNumber,
		Name:      *res.Payload[0].Stop.Name,
		Lat:       float32(*res.Payload[0].Stop.Location.Latitude),
		Lon:       float32(*res.Payload[0].Stop.Location.Longitude),
	})
	p.parseDepartureArrival(res.Payload, station.EvaNumber, true)
}

func (p *DbRest) requestDeparture(station providers.ProviderStation, when time.Time, duration int64) {
	var params = operations.NewGetStopsIDDeparturesParams()
	params.ID = strconv.Itoa(station.EvaNumber)
	params.Duration = &duration
	params.When = (*strfmt.DateTime)(&when)

	res, err := p.client.Operations.GetStopsIDDepartures(params)
	if err != nil {
		log.Panic(err)
		return
	}
	p.parseDepartureArrival(res.Payload, station.EvaNumber, false)
}

func (p *DbRest) parseDepartureArrival(stops []*models.DepartureArrival, evaNumber int, arrival bool) {
	for _, stop := range stops {
		lineID, err := strconv.Atoi(*stop.Line.FahrtNr)
		if err != nil {
			log.Printf("Failed to convert Line ID %d", stop.Line.FahrtNr)
		}
		p.parseLine(stop, lineID)
		p.parseLineStop(stop, arrival, evaNumber, lineID)
	}
}

func (p *DbRest) parseLine(stop *models.DepartureArrival, lineID int) {
	lineName := ""
	if stop.Line.Name != nil {
		lineName = *stop.Line.Name
	}
	productName := ""
	if stop.Line.ProductName != nil {
		productName = *stop.Line.ProductName
	}
	p.consumer.UpsertLine(providers.ProviderLine{ID: lineID, Type: productName, Name: lineName})
}

func (p *DbRest) parseLineStop(stop *models.DepartureArrival, arrival bool, evaNumber int, lineID int) {

	planned := &providers.ProviderLineStopInfo{}
	current := &providers.ProviderLineStopInfo{}

	if arrival {
		if stop.PlannedWhen != nil {
			planned.Arrival = time.Time(*stop.PlannedWhen)
		}
		if stop.PlannedPlatform != nil {
			planned.ArrivalTrack = *stop.PlannedPlatform
		}
		if stop.When != nil {
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
		if stop.When != nil {
			current.Departure = time.Time(*stop.When)
		}
		if stop.Platform != nil {
			current.DepartureTrack = *stop.Platform
		}
	}
	p.consumer.UpsertLineStop(providers.ProviderLineStop{EvaNumber: evaNumber, LineID: lineID, Planned: planned, Current: current})
}

func (p *DbRest) prepareClient() {
	r := httptransport.New(os.Getenv("API_CACHE_HOST"), os.Getenv("DBREST_API_CACHE_PREFIX"), []string{os.Getenv("API_CACHE_SCHEME")})
	r.DefaultAuthentication = httptransport.BearerToken(os.Getenv("DB_API_ACCESS_TOKEN"))
	p.client = apiclient.New(r, strfmt.Default)
}
