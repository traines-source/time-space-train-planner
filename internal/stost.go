package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	proto "google.golang.org/protobuf/proto"
	"traines.eu/time-space-train-planner/internal/stost"
)

var productTypes = map[string]int32{
	"bus":             10,
	"subway":          7,
	"tram":            9,
	"regional":        6,
	"suburban":        8,
	"nationalExpress": 1,
	"national":        2,
	"regionalExpress": 5,
	"ferry":           11,
	"taxi":            12,
	"regionalExp":     5,
	"Foot":            100,
}

func toDelay(planned time.Time, current time.Time) *int32 {
	if current.IsZero() {
		return nil
	}
	d := int32(math.Round(current.Sub(planned).Minutes()))
	return &d
}

func toCurrent(stopInfo *stost.StopInfo, actual bool) time.Time {
	if !stopInfo.IsLive && !actual {
		return time.Time{}
	}
	return time.Unix(stopInfo.Scheduled, 0).Add(time.Duration(stopInfo.GetDelayMinutes()) * time.Minute)
}

func createRequestMessage(system string, from string, to string, startTime time.Time, now time.Time) *stost.Message {
	//d, _ := time.ParseDuration("0h")
	return &stost.Message{
		Query: &stost.Query{
			Origin:      from,
			Destination: to,
			Now:         now.Unix(), //startTime.Add(d).Unix(),
		},
		Timetable: &stost.Timetable{
			Stations:  []*stost.Station{},
			Routes:    []*stost.Route{},
			StartTime: startTime.Unix(),
		},
		System: system,
	}
}

func prepareForEnrichment(requestMessage *stost.Message, lines map[string]*Line, stations map[string]*Station, regionly bool) {
	for _, s := range stations {
		requestMessage.Timetable.Stations = append(requestMessage.Timetable.Stations, &stost.Station{
			Id:   s.ID,
			Name: s.Name,
			Lat:  &s.Lat,
			Lon:  &s.Lon,
		})
	}
	for _, l := range lines {
		var p int32 = -1
		if l.Type == "Foot" && len(l.Route) > 0 && l.Route[0].Discarded {
			continue
		}
		if regionly && (l.Type == "national" || l.Type == "nationalExpress") {
			continue
		}
		if val, ok := productTypes[l.Type]; ok {
			p = val
		}
		route := &stost.Route{
			Id:          l.ID,
			Name:        l.Name,
			ProductType: p,
			Direction:   &l.Direction,
			Message:     &l.Message,
			Trips: []*stost.Trip{{
				Connections: []*stost.Connection{},
			}},
		}
		for _, c := range l.Route {
			route.Trips[0].Connections = append(route.Trips[0].Connections, &stost.Connection{
				FromId:    c.From.ID,
				ToId:      c.To.ID,
				Cancelled: c.Cancelled,
				Message:   &c.Message,
				Departure: &stost.StopInfo{
					Scheduled:    c.Planned.Departure.Unix(),
					DelayMinutes: toDelay(c.Planned.Departure, c.Current.Departure),
					IsLive:       !c.Current.Departure.IsZero(),
				},
				Arrival: &stost.StopInfo{
					Scheduled:    c.Planned.Arrival.Unix(),
					DelayMinutes: toDelay(c.Planned.Arrival, c.Current.Arrival),
					IsLive:       !c.Current.Arrival.IsZero(),
				},
			})
		}
		requestMessage.Timetable.Routes = append(requestMessage.Timetable.Routes, route)
	}
}

func queryStost(requestMessage *stost.Message) *stost.Message {
	out, err := proto.Marshal(requestMessage)
	if err != nil {
		log.Println("Failed to encode proto message:", err)
		return nil
	}

	url := os.Getenv("STOST_API_URL")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(out))
	if err != nil {
		log.Println(err)
		return nil
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	responseMessage := &stost.Message{}
	if err := proto.Unmarshal(body, responseMessage); err != nil {
		log.Println("Failed to parse responseMessage:", err)
		return nil
	}
	return responseMessage
}

func convertDistribution(distr *stost.Distribution) Distribution {
	if distr == nil {
		return Distribution{}
	}
	d := Distribution{
		Histogram:           distr.Histogram,
		Start:               time.Unix(distr.Start, 0),
		FeasibleProbability: distr.FeasibleProbability,
	}
	if len(distr.Histogram) > 0 {
		d.Mean = time.Unix(distr.Mean, 0)
	}
	return d
}

func enrichWithDistributions(responseMessage *stost.Message, lines map[string]*Line) {
	for _, r := range responseMessage.Timetable.Routes {
		for _, t := range r.Trips {
			for i, c := range t.Connections {
				lines[r.Id].Route[i].DestinationArrival = convertDistribution(c.DestinationArrival)
			}
		}
	}
}

func reverseMap(m map[string]int32) map[int32]string {
	newM := map[int32]string{}
	for k, v := range m {
		newM[v] = k
	}
	return newM
}

func produce(responseMessage *stost.Message, lines map[string]*Line, stations map[string]*Station) {
	for _, s := range responseMessage.Timetable.Stations {
		var group *string = nil
		if s.Parent != nil && *s.Parent != "" {
			group = s.Parent
		}
		stations[s.Id] = &Station{
			ID:      s.Id,
			Name:    s.Name,
			Lat:     s.GetLat(),
			Lon:     s.GetLon(),
			GroupID: group,
		}
	}
	productTypesReverse := reverseMap(productTypes)
	i := 0
	for _, r := range responseMessage.Timetable.Routes {
		for _, t := range r.Trips {
			line := Line{
				ID:        strconv.Itoa(i),
				Name:      r.Name,
				Type:      productTypesReverse[r.ProductType],
				Direction: r.GetDirection(),
				Message:   r.GetMessage(),
			}
			lines[line.ID] = &line
			for _, c := range t.Connections {
				from := stations[c.FromId]
				to := stations[c.ToId]
				distr := convertDistribution(c.DestinationArrival)
				edge := &Edge{
					Line:    &line,
					From:    from,
					To:      to,
					Message: c.GetMessage(),
					Planned: StopInfo{
						Departure:      time.Unix(c.Departure.Scheduled, 0),
						DepartureTrack: c.Departure.GetScheduledTrack(),
						Arrival:        time.Unix(c.Arrival.Scheduled, 0),
						ArrivalTrack:   c.Arrival.GetScheduledTrack(),
					},
					Current: StopInfo{
						Departure:      toCurrent(c.Departure, false),
						DepartureTrack: c.Departure.GetScheduledTrack(),
						Arrival:        toCurrent(c.Arrival, false),
						ArrivalTrack:   c.Arrival.GetScheduledTrack(),
					},
					Actual: StopInfo{
						Departure:      toCurrent(c.Departure, true),
						DepartureTrack: c.Departure.GetScheduledTrack(),
						Arrival:        toCurrent(c.Arrival, true),
						ArrivalTrack:   c.Arrival.GetScheduledTrack(),
					},
					Cancelled:                  c.Cancelled,
					DestinationArrival:         distr,
					EarliestDestinationArrival: distr.Mean,
				}
				line.Route = append(line.Route, edge)
				from.Departures = append(from.Departures, edge)
				to.Arrivals = append(to.Arrivals, edge)
			}
			i++
		}
	}
}

func StostEnrich(system string, lines map[string]*Line, stations map[string]*Station, from string, to string, startTime time.Time, now time.Time, regionly bool) {
	requestMessage := createRequestMessage("de_db", from, to, startTime, now)
	prepareForEnrichment(requestMessage, lines, stations, regionly)
	responseMessage := queryStost(requestMessage)
	if responseMessage == nil {
		return
	}
	enrichWithDistributions(responseMessage, lines)
}

func StostProduce(system string, lines map[string]*Line, stations map[string]*Station, from string, to string, startTime time.Time, now time.Time) {
	log.Print("Producing based on stost...")
	requestMessage := createRequestMessage(system, from, to, startTime, now)
	responseMessage := queryStost(requestMessage)
	if responseMessage == nil {
		log.Panic("Empty response, stopping.")
	}
	produce(responseMessage, lines, stations)
}
