package dbrest

import (
	"log"
	"os"
	"strconv"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"traines.eu/time-space-train-planner/providers"
	apiclient "traines.eu/time-space-train-planner/providers/dbrest/client"
	"traines.eu/time-space-train-planner/providers/dbrest/client/operations"
)

type DbRest struct {
	consumer providers.Consumer
	client   *apiclient.Dbrest
}

func (p *DbRest) Fetch(c providers.Consumer) {
	p.consumer = c
	p.prepareClient()
	p.requestStations()
}

func (p *DbRest) requestStations() {
	stations := p.consumer.Stations()
	for _, station := range stations {
		p.requestStation(station)
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

func (p *DbRest) prepareClient() {
	r := httptransport.New(os.Getenv("API_CACHE_HOST"), os.Getenv("API_CACHE_PREFIX")+"/dbrest", apiclient.DefaultSchemes)
	r.DefaultAuthentication = httptransport.BearerToken(os.Getenv("DB_API_ACCESS_TOKEN"))
	p.client = apiclient.New(r, strfmt.Default)
}
