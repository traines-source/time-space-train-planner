package render

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"text/template"

	"traines.eu/time-space-train-planner/internal"
)

const MAX_VIAS = 10

type model struct {
	From      StationLabel
	To        StationLabel
	Vias      []StationLabel
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

func Vias(stations map[int]*internal.Station, from int, to int, dateTime string, wr io.Writer, err error) {
	m := &model{
		LegalLink: os.Getenv("TSTP_LEGAL"),
		From:      makeStationLabel(stations[from]),
		To:        makeStationLabel(stations[to]),
		Vias:      makeVias(stations, from, to),
		DateTime:  dateTime,
		Error:     err,
	}
	var l = len(m.Vias)
	fillupStations(m, l)
	m.template(wr)
}

func ViasApi(stations map[int]*internal.Station, from int, to int, dateTime string, wr http.ResponseWriter, err error) {
	m := &model{
		From:     makeStationLabel(stations[from]),
		To:       makeStationLabel(stations[to]),
		Vias:     makeVias(stations, from, to),
		DateTime: dateTime,
		Error:    err,
	}
	json.NewEncoder(wr).Encode(m)
}

func fillupStations(m *model, existing int) {
	for i := 0; i < 10-existing; i++ {
		m.Vias = append(m.Vias, StationLabel{})
	}
}

func (m *model) template(wr io.Writer) {
	t := template.Must(template.New("index.tmpl.html").ParseFiles("./render/index.tmpl.html"))
	err := t.Execute(wr, m)
	if err != nil {
		panic(err)
	}
}
