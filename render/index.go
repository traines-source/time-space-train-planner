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
}

func Index(wr io.Writer) {
	m := &model{}
	m.template(wr)
}

func Vias(stations map[int]*internal.Station, from int, to int, wr io.Writer) {
	m := &model{
		From: stations[from],
		To:   stations[to],
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
	for i := 0; i < 10-l; i++ {
		m.Stations = append(m.Stations, &internal.Station{})
	}
	m.template(wr)
}

func (m *model) template(wr io.Writer) {
	t := template.Must(template.New("index.tmpl.html").ParseFiles("./render/index.tmpl.html"))
	err := t.Execute(wr, m)
	if err != nil {
		panic(err)
	}
}
