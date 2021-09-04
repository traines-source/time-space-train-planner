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
	Edges            map[[2]*internal.Edge]*EdgePath
	minTime          time.Time
	maxTime          time.Time
	timeAxisDistance float32
	TimeIndicators   []time.Time
	TimeAxisSize     int
	SpaceAxisSize    int
}

const (
	timeAxisSize             = 2000
	spaceAxisSize            = 2000
	maxTimeIndicators        = 5
	minTimeIndicatorDistance = "15m"
)

func TimeSpace(stations map[int]*internal.Station, lines map[int]*internal.Line, wr io.Writer) {
	c := &container{Stations: map[*internal.Station]*StationLabel{}, Edges: map[[2]*internal.Edge]*EdgePath{}}
	c.setupStations(stations)
	c.setupEdges(lines)
	c.gravitate()
	c.render(wr)
}

func (c *container) setupStations(stations map[int]*internal.Station) {
	for _, s := range stations {
		station := &StationLabel{Station: *s}
		station.Coord.SpaceAxis = station
		c.Stations[s] = station

	}
}

func (c *container) setupEdges(lines map[int]*internal.Line) {
	log.Print("test")
	for _, l := range lines {
		for i := 0; i < len(l.Route); i++ {
			e := l.Route[i]
			log.Print("e", e.ShortestPath)
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
	for _, l := range lines {
		for i := 0; i < len(l.Route); i++ {
			origin := l.Route[i]
			if originEdgePath, ok := c.Edges[[2]*internal.Edge{origin, origin}]; ok {
				var lastEdge *internal.Edge
				for e := origin.ShortestPath; e != nil; e = e.ShortestPath {
					if edgePath, ok := c.Edges[[2]*internal.Edge{e, e}]; ok {
						edgePath.ShortestPathFor = append(edgePath.ShortestPathFor, originEdgePath)
						log.Print("here", len(edgePath.ShortestPathFor))

					}
					if edgePath, ok := c.Edges[[2]*internal.Edge{lastEdge, e}]; ok {
						edgePath.ShortestPathFor = append(edgePath.ShortestPathFor, originEdgePath)
					}
					lastEdge = e
				}
			}
		}
	}
	//log.Print("test")
	//log.Printf("%+v", lines)
	//log.Printf("%+v", c.Edges)
}

func (c *container) stretchTimeAxis(min time.Time, max time.Time) {
	if min.Before(c.minTime) || c.minTime.IsZero() {
		c.minTime = min
	}
	if max.After(c.maxTime) || c.maxTime.IsZero() {
		c.maxTime = max
	}
}

func (c *container) insertEdge(e *internal.Edge) *EdgePath {
	edge := &EdgePath{
		Edge: *e,
		ID:   fmt.Sprintf("%d_%d", e.Line.ID, e.From.EvaNumber),
		From: Coord{SpaceAxis: c.Stations[e.From], TimeAxis: e.Actual.Departure},
		To:   Coord{SpaceAxis: c.Stations[e.To], TimeAxis: e.Actual.Arrival},
	}
	c.Edges[[2]*internal.Edge{e, e}] = edge
	return edge
}

func (c *container) insertStationEdge(last *internal.Edge, this *internal.Edge) *EdgePath {
	edge := &EdgePath{
		From: Coord{SpaceAxis: c.Stations[this.From], TimeAxis: last.Actual.Arrival},
		To:   Coord{SpaceAxis: c.Stations[this.From], TimeAxis: this.Actual.Departure},
	}
	c.Edges[[2]*internal.Edge{last, this}] = edge
	return edge
}

func (c *container) gravitate() {
	c.TimeAxisSize = timeAxisSize
	c.SpaceAxisSize = spaceAxisSize
	num := float32(len(c.Stations))
	for _, s := range c.Stations {
		s.SpaceAxis = int(float32(s.Station.Rank)/num*float32(c.SpaceAxisSize-50) + 50.0)
	}
	c.indicateTimes()
}

func (c *container) indicateTimes() {
	delta := c.maxTime.Unix() - c.minTime.Unix()
	c.timeAxisDistance = float32(delta)
	//duration, _ := time.ParseDuration(fmt.Sprintf("%ds", delta/maxTimeIndicators))
	//now := c.minTime
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
		return 50
	}
	delta := float32(t.Unix() - c.minTime.Unix())
	//log.Print(t, delta, c.timeAxisDistance)
	return int(delta/c.timeAxisDistance*float32(c.TimeAxisSize-100) + 100.0)
}

func (c *container) X(coord Coord) int {
	return coord.SpaceAxis.SpaceAxis

}

func (c *container) Y(coord Coord) int {
	return c.timeAxis(coord.TimeAxis)
}

func (c *container) Label(p *EdgePath) string {
	e := p.Edge
	if e.Line == nil {
		return ""
	}
	var label string
	if e.Line.Name != "" {
		label = e.Line.Name
	} else {
		label = fmt.Sprintf("%d", e.Line.ID)
	}
	return e.Line.Type + " " + label
}

func (c *container) Departure(p *EdgePath) string {
	e := p.Edge
	if e.Line == nil {
		return ""
	}
	var label string
	label = simpleTime(e.Actual.Departure)
	if e.Planned.DepartureTrack != "" {
		label += e.Planned.DepartureTrack
	}
	return label
}

func (c *container) Arrival(p *EdgePath) string {
	e := p.Edge
	if e.Line == nil {
		return ""
	}
	label := simpleTime(e.Actual.Arrival)
	return label
}

func simpleTime(t time.Time) string {
	return fmt.Sprintf("%02d:%02d ", t.Hour(), t.Minute())
}
