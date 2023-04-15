package render

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"text/template"

	"traines.eu/time-space-train-planner/internal"
)

const MAX_VIAS = 10

type model struct {
	From      StationLabel
	To        StationLabel
	Stations  []StationLabel
	DateTime  string
	LegalLink string
	Error     error
}

func Index(wr io.Writer) {
	m := &model{
		LegalLink: os.Getenv("TSTP_LEGAL"),
		From:      StationLabel{},
		To:        StationLabel{},
	}
	m.template(wr)
}

func makeStationLabel(s *internal.Station) StationLabel {
	return StationLabel{ID: strconv.Itoa(s.EvaNumber), Name: s.Name, Rank: s.Rank}
}

func Vias(stations map[int]*internal.Station, from int, to int, dateTime string, wr io.Writer, err error) {
	m := &model{
		LegalLink: os.Getenv("TSTP_LEGAL"),
		From:      makeStationLabel(stations[from]),
		To:        makeStationLabel(stations[to]),
		DateTime:  dateTime,
		Error:     err,
	}
	populateStations(stations, from, to, m)
	var l = len(m.Stations)
	fillupStations(m, l)
	m.template(wr)
}

func ViasApi(stations map[int]*internal.Station, from int, to int, dateTime string, wr http.ResponseWriter, err error) {
	m := &model{
		From:     makeStationLabel(stations[from]),
		To:       makeStationLabel(stations[to]),
		Stations: []StationLabel{},
		DateTime: dateTime,
		Error:    err,
	}
	populateStations(stations, from, to, m)
	json.NewEncoder(wr).Encode(m)
}

func populateStations(stations map[int]*internal.Station, from int, to int, m *model) {
	for _, s := range stations {
		log.Print(s)
		if s.EvaNumber == from || s.EvaNumber == to || s.GroupNumber != nil && *s.GroupNumber != s.EvaNumber {
			continue
		}
		m.Stations = append(m.Stations, makeStationLabel(s))
	}
	sort.Slice(m.Stations, func(i, j int) bool {
		return m.Stations[i].Rank < m.Stations[j].Rank
	})
}

func fillupStations(m *model, existing int) {
	for i := 0; i < 10-existing; i++ {
		m.Stations = append(m.Stations, StationLabel{})
	}
}

func (m *model) template(wr io.Writer) {
	t := template.Must(template.New("index.tmpl.html").ParseFiles("./render/index.tmpl.html"))
	err := t.Execute(wr, m)
	if err != nil {
		panic(err)
	}
}
