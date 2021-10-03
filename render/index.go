package render

import (
	"io"
	"sort"
	"text/template"

	"traines.eu/time-space-train-planner/internal"
)

const MAX_VIAS = 10

type model struct {
	From     *internal.Station
	To       *internal.Station
	Stations []*internal.Station
	DateTime string
}

func Index(wr io.Writer) {
	m := &model{
		From: &internal.Station{},
		To:   &internal.Station{},
	}
	m.template(wr)
}

func Vias(stations map[int]*internal.Station, from int, to int, dateTime string, wr io.Writer) {
	m := &model{
		From:     stations[from],
		To:       stations[to],
		DateTime: dateTime,
	}
	for _, s := range stations {
		if s.EvaNumber == from || s.EvaNumber == to {
			continue
		}
		m.Stations = append(m.Stations, s)
	}
	sort.Slice(m.Stations, func(i, j int) bool {
		return m.Stations[i].Rank < m.Stations[j].Rank
	})
	var l = len(m.Stations)
	fillupStations(m, l)
	m.template(wr)
}

func fillupStations(m *model, existing int) {
	for i := 0; i < 10-existing; i++ {
		m.Stations = append(m.Stations, &internal.Station{})
	}
}

func (m *model) template(wr io.Writer) {
	t := template.Must(template.New("index.tmpl.html").ParseFiles("./render/index.tmpl.html"))
	err := t.Execute(wr, m)
	if err != nil {
		panic(err)
	}
}
