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
}

func (p *Timetables) Fetch(c providers.Consumer) {
	client := prepareClient()
	stations := c.Stations()
	for _, station := range stations {
		var params = operations.NewGetPlanEvaNoDateHourParams()
		params.EvaNo = strconv.Itoa(station.EvaNumber)
		t := time.Now()
		params.Date = fmt.Sprintf("%02d%02d%02d", t.Year()%100, t.Month(), t.Day())
		params.Hour = fmt.Sprintf("%02d", t.Hour())
		res, err := client.Operations.GetPlanEvaNoDateHour(params)
		if err != nil {
			log.Panic(err)
			return
		}
		c.UpsertStation(providers.ProviderStation{EvaNumber: station.EvaNumber, Name: res.Payload.Station})

		for _, stop := range res.Payload.S {
			lineID, err := strconv.Atoi(*stop.Tl.N)
			if err != nil {
				log.Printf("Failed to convert Line ID %d", lineID)
			}
			p.parseLine(stop, lineID, c)
			p.parseLineStop(stop, station.EvaNumber, lineID, c)
		}
	}
}

func (p *Timetables) parseLine(stop *models.TimetableStop, lineID int, c providers.Consumer) {
	lineName := ""
	if stop.Ar != nil {
		lineName = stop.Ar.L
	}
	if stop.Dp != nil {
		lineName = stop.Dp.L
	}
	c.UpsertLine(providers.ProviderLine{ID: lineID, Type: *stop.Tl.C, Name: lineName})
}

func (p *Timetables) parseLineStop(stop *models.TimetableStop, evaNumber int, lineID int, c providers.Consumer) {

	planned := &providers.ProviderLineStopInfo{}
	if stop.Ar != nil {
		planned.Arrival = parseEventTime(stop.Ar.Pt)
		planned.Track = stop.Ar.Pp
	}
	if stop.Dp != nil {
		planned.Departure = parseEventTime(stop.Dp.Pt)
		planned.Track = stop.Dp.Pp
	}
	c.UpsertLineStop(providers.ProviderLineStop{EvaNumber: evaNumber, LineID: lineID, Planned: planned})
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

func prepareClient() *apiclient.Timetable {
	r := httptransport.New(apiclient.DefaultHost, apiclient.DefaultBasePath, apiclient.DefaultSchemes)
	r.DefaultAuthentication = httptransport.BearerToken(os.Getenv("DB_API_ACCESS_TOKEN"))
	r.DefaultMediaType = runtime.XMLMime
	r.Consumers = map[string]runtime.Consumer{
		runtime.XMLMime: runtime.XMLConsumer(),
	}
	r.Producers = map[string]runtime.Producer{
		"application/xhtml+xml": runtime.XMLProducer(),
	}
	return apiclient.New(r, strfmt.Default)
}
