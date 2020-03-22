package render

import (
	"fmt"
	"io"
	"log"
	"text/template"
	"time"

	"traines.eu/time-space-train-planner/internal"
)

type container struct {
	Stations         map[*internal.Station]*StationLabel
	Edges            []*EdgePath
	minTime          time.Time
	maxTime          time.Time
	timeAxisDistance float32
	Width            int
	Height           int
}

const (
	width  = 500
	height = 500
)

func TimeSpace(stations map[int]*internal.Station, lines map[int]*internal.Line, wr io.Writer) {
	c := &container{}
	c.setupStations(stations)
	c.setupEdges(lines)
	c.gravitate()
	c.render(wr)
}

func (c *container) setupStations(stations map[int]*internal.Station) {
	c.Stations = map[*internal.Station]*StationLabel{}
	for _, s := range stations {
		station := &StationLabel{Station: *s}
		station.Coord.SpaceAxis = station
		c.Stations[s] = station

	}
}

func (c *container) setupEdges(lines map[int]*internal.Line) {
	for _, l := range lines {
		for i := 0; i < len(l.Route); i++ {
			e := l.Route[i]
			log.Print("e", e)
			if e.Redundant {
				continue
			}
			if i > 0 {
				c.insertStationEdge(l.Route[i-1], e)
			}
			edge := c.insertEdge(e)
			c.stretchTimeAxis(edge.From.TimeAxis, edge.To.TimeAxis)
		}
	}
	log.Print("test")
	log.Print("%+v", lines)
	log.Print("%+v", c.Edges)
}

func (c *container) stretchTimeAxis(min time.Time, max time.Time) {
	if min.Before(c.minTime) {
		c.minTime = min
	}
	if max.After(c.maxTime) {
		c.maxTime = max
	}
}

func (c *container) insertEdge(e *internal.Edge) *EdgePath {
	edge := &EdgePath{
		Edge:  *e,
		ID:    fmt.Sprintf("%d_%d", e.Line.ID, e.From.EvaNumber),
		From:  Coord{SpaceAxis: c.Stations[e.From], TimeAxis: e.Actual.Departure},
		To:    Coord{SpaceAxis: c.Stations[e.To], TimeAxis: e.Actual.Arrival},
		Label: makeLabel(e),
	}
	c.Edges = append(c.Edges, edge)
	return edge
}

func (c *container) insertStationEdge(last *internal.Edge, this *internal.Edge) *EdgePath {
	edge := &EdgePath{
		From: Coord{SpaceAxis: c.Stations[this.From], TimeAxis: last.Actual.Arrival},
		To:   Coord{SpaceAxis: c.Stations[this.From], TimeAxis: this.Actual.Departure},
	}
	c.Edges = append(c.Edges, edge)
	return edge
}

func makeLabel(e *internal.Edge) string {
	var label string
	if e.Line.Name != "" {
		label = e.Line.Name
	} else {
		label = fmt.Sprintf("%d", e.Line.ID)
	}
	return e.Line.Type + " " + label
}

func (c *container) gravitate() {
	c.Width = width
	c.Height = height
	num := float32(len(c.Stations))
	i := 0
	for _, s := range c.Stations {
		s.SpaceAxis = int(float32(i)/num*float32(c.Height-50) + 50.0)
		i++
	}
	c.timeAxisDistance = float32(c.maxTime.Unix() - c.minTime.Unix())
}

func (c *container) render(wr io.Writer) {
	t := template.Must(template.New("time-space.tmpl").ParseFiles("./render/time-space.tmpl"))
	err := t.Execute(wr, c)
	if err != nil {
		panic(err)
	}
}

func (c *container) timeAxis(t time.Time) int {
	if t.IsZero() {
		return 0
	}
	return int(float32((t.Unix()-c.minTime.Unix()))/c.timeAxisDistance*float32(c.Width-100) + 100.0)
}

func (c *container) X(coord Coord) int {
	return c.timeAxis(coord.TimeAxis)
}

func (c *container) Y(coord Coord) int {
	return coord.SpaceAxis.SpaceAxis
}
