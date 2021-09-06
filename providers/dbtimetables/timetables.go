package dbtimetables

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"traines.eu/time-space-train-planner/providers"
	apiclient "traines.eu/time-space-train-planner/providers/dbtimetables/client"
	"traines.eu/time-space-train-planner/providers/dbtimetables/client/operations"
	"traines.eu/time-space-train-planner/providers/dbtimetables/models"
)

type Timetables struct {
	consumer providers.Consumer
	client   *apiclient.Timetable
}

func (p *Timetables) Fetch(c providers.Consumer) {
	p.consumer = c
	p.prepareClient()

	current, to := p.consumer.RequestStationDataBetween(nil)
	delta, _ := time.ParseDuration("1h")

	i := 0
	for current.Before(to) {
		p.requestAtTime(current)
		//return
		current = current.Add(delta)
		i++
		if i > 3 {
			break
		}
	}
}

func (p *Timetables) requestAtTime(time time.Time) {
	stations := p.consumer.Stations()
	for _, station := range stations {
		from, to := p.consumer.RequestStationDataBetween(&station)
		if from.Equal(time) || from.Before(time) && time.Before(to) {
			p.requestStationAtTime(station, time)
		}
	}
}

func (p *Timetables) requestStationAtTime(station providers.ProviderStation, t time.Time) {
	var params = operations.NewGetPlanEvaNoDateHourParams()
	params.EvaNo = strconv.Itoa(station.EvaNumber)
	params.Date = fmt.Sprintf("%02d%02d%02d", t.Year()%100, t.Month(), t.Day())
	params.Hour = fmt.Sprintf("%02d", t.Hour())
	res, err := p.client.Operations.GetPlanEvaNoDateHour(params)
	if err != nil {
		log.Panic(err)
		return
	}
	p.consumer.UpsertStation(providers.ProviderStation{EvaNumber: station.EvaNumber, Name: res.Payload.Station})

	for _, stop := range res.Payload.S {
		lineID, err := strconv.Atoi(*stop.Tl.N)
		if err != nil {
			log.Printf("Failed to convert Line ID %d", lineID)
		}
		p.parseLine(stop, lineID)
		p.parseLineStop(stop, station.EvaNumber, lineID)
	}
}

func (p *Timetables) parseLine(stop *models.TimetableStop, lineID int) {
	lineName := ""
	if stop.Ar != nil {
		lineName = stop.Ar.L
	}
	if stop.Dp != nil {
		lineName = stop.Dp.L
	}
	p.consumer.UpsertLine(providers.ProviderLine{ID: lineID, Type: *stop.Tl.C, Name: lineName})
}

func (p *Timetables) parseLineStop(stop *models.TimetableStop, evaNumber int, lineID int) {

	planned := &providers.ProviderLineStopInfo{}
	if stop.Ar != nil {
		planned.Arrival = parseEventTime(stop.Ar.Pt)
		planned.ArrivalTrack = stop.Ar.Pp
	}
	if stop.Dp != nil {
		planned.Departure = parseEventTime(stop.Dp.Pt)
		planned.DepartureTrack = stop.Dp.Pp
	}
	p.consumer.UpsertLineStop(providers.ProviderLineStop{EvaNumber: evaNumber, LineID: lineID, Planned: planned})
}

func parseEventTime(timeString string) time.Time {
	r := []rune(timeString)
	century := time.Now().Year() / 100 * 100
	year := century + ato2i(r, 0)
	return time.Date(year, time.Month(ato2i(r, 2)), ato2i(r, 4), ato2i(r, 6), ato2i(r, 8), 0, 0, time.Local)
}

func ato2i(r []rune, index int) int {
	val, err := strconv.Atoi(string(r[index : index+2]))
	if err != nil {
		log.Panicf("failed to parse date %s", string(r))
	}
	return val
}

func (p *Timetables) prepareClient() {
	r := httptransport.New(os.Getenv("API_CACHE_HOST"), os.Getenv("API_CACHE_PREFIX")+apiclient.DefaultBasePath, apiclient.DefaultSchemes)
	r.DefaultAuthentication = httptransport.BearerToken(os.Getenv("DB_API_ACCESS_TOKEN"))
	r.DefaultMediaType = runtime.XMLMime
	r.Consumers = map[string]runtime.Consumer{
		runtime.XMLMime: runtime.XMLConsumer(),
	}
	r.Producers = map[string]runtime.Producer{
		"application/xhtml+xml": runtime.XMLProducer(),
	}
	p.client = apiclient.New(r, strfmt.Default)
}
