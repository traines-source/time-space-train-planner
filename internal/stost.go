package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
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

func delay(planned time.Time, current time.Time) *int32 {
	if current.IsZero() {
		return nil
	}
	d := int32(math.Round(current.Sub(planned).Minutes()))
	return &d
}

func StostEnrich(lines map[string]*Line, stations map[string]*Station, from string, to string, startTime time.Time, now time.Time) {
	//d, _ := time.ParseDuration("0h")
	requestMessage := &stost.Message{
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
		System: "de_db",
	}
	for _, s := range stations {
		requestMessage.Timetable.Stations = append(requestMessage.Timetable.Stations, &stost.Station{
			Id:   s.ID,
			Name: s.Name,
			Lat: &s.Lat,
			Lon: &s.Lon,
		})
	}
	for _, l := range lines {
		var p int32 = -1
		if l.Type == "Foot" {
			continue
		}
		if val, ok := productTypes[l.Type]; ok {
			p = val
		}
		route := &stost.Route{
			Id:          l.ID,
			Name:        l.Name,
			ProductType: p,
			Trips: []*stost.Trip{{
				Connections: []*stost.Connection{},
			}},
		}
		for _, c := range l.Route {
			if c.Discarded {
				continue
			}
			route.Trips[0].Connections = append(route.Trips[0].Connections, &stost.Connection{
				FromId:    c.From.ID,
				ToId:      c.To.ID,
				Cancelled: c.Cancelled,
				Departure: &stost.StopInfo{
					Scheduled:    c.Planned.Departure.Unix(),
					DelayMinutes: delay(c.Planned.Departure, c.Current.Departure),
					IsLive:       !c.Current.Departure.IsZero(),
				},
				Arrival: &stost.StopInfo{
					Scheduled:    c.Planned.Arrival.Unix(),
					DelayMinutes: delay(c.Planned.Arrival, c.Current.Arrival),
					IsLive:       !c.Current.Arrival.IsZero(),
				},
			})
		}
		requestMessage.Timetable.Routes = append(requestMessage.Timetable.Routes, route)
	}

	out, err := proto.Marshal(requestMessage)
	if err != nil {
		log.Println("Failed to encode proto message:", err)
		return
	}

	url := os.Getenv("STOST_API_URL")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(out))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	responseMessage := &stost.Message{}
	if err := proto.Unmarshal(body, responseMessage); err != nil {
		log.Println("Failed to parse responseMessage:", err)
		return
	}
	for _, r := range responseMessage.Timetable.Routes {
		for _, t := range r.Trips {
			for i, c := range t.Connections {
				if c.DestinationArrival == nil {
					continue
				}
				lines[r.Id].Route[i].DestinationArrival = Distribution{
					Histogram:           c.DestinationArrival.Histogram,
					Start:               time.Unix(c.DestinationArrival.Start, 0),
					Mean:                time.Unix(c.DestinationArrival.Mean, 0),
					FeasibleProbability: c.DestinationArrival.FeasibleProbability,
				}
			}
		}
	}
}
