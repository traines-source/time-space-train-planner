package render

import (
	"fmt"
	"io"
	"log"
	"sort"
	"text/template"
	"time"

	"traines.eu/time-space-train-planner/internal"
)

type container struct {
	Stations         map[*internal.Station]*StationLabel
	Edges            map[string]*EdgePath
	SortedEdges      []*EdgePath
	minTime          time.Time
	maxTime          time.Time
	maxSpace         int
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

func TimeSpace(stations map[int]*internal.Station, lines map[string]*internal.Line, wr io.Writer) {
	c := &container{Stations: map[*internal.Station]*StationLabel{}, Edges: map[string]*EdgePath{}}
	c.setupStations(stations)
	c.setupEdges(lines)
	c.gravitate()
	c.render(wr)
}

func (c *container) setupStations(stations map[int]*internal.Station) {
	for _, s := range stations {
		if len(s.Arrivals) > 0 || len(s.Departures) > 0 {
			station := &StationLabel{Station: *s}
			station.Coord.SpaceAxis = station
			c.Stations[s] = station
		}
	}
}

func (c *container) setupEdges(lines map[string]*internal.Line) {
	for _, l := range lines {
		for i := 0; i < len(l.Route); i++ {
			e := l.Route[i]
			if i > 0 {
				c.insertStationEdge(l.Route[i-1], e)
			}
			edge := c.insertEdge(e)
			c.stretchTimeAxis(edge.From.TimeAxis, edge.To.TimeAxis)
		}
	}
	sort.Slice(c.SortedEdges, func(i, j int) bool {
		if c.SortedEdges[i].Redundant == c.SortedEdges[j].Redundant && c.SortedEdges[i].Line != nil {
			return c.SortedEdges[i].Line.Type == "Foot"
		}
		return c.SortedEdges[i].Redundant
	})
	for _, l := range lines {
		for i := 0; i < len(l.Route); i++ {
			origin := l.Route[i]
			if originEdgePath, ok := c.Edges[generateEdgeID(origin)]; ok {
				var lastEdge *internal.Edge
				for e := origin; e != nil; e = e.ShortestPath {
					if edgePath, ok := c.Edges[generateEdgeID(e)]; ok {
						edgePath.ShortestPathFor = append(edgePath.ShortestPathFor, originEdgePath)
					}
					if edgePath, ok := c.Edges[generateStationEdgeID(lastEdge, e)]; ok {

						edgePath.ShortestPathFor = append(edgePath.ShortestPathFor, originEdgePath)
					}
					lastEdge = e
				}
			}
		}
	}
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
		ID:   generateEdgeID(e),
		From: Coord{SpaceAxis: c.Stations[e.From], TimeAxis: e.Actual.Departure},
		To:   Coord{SpaceAxis: c.Stations[e.To], TimeAxis: e.Actual.Arrival},
	}
	c.Edges[edge.ID] = edge
	c.SortedEdges = append(c.SortedEdges, edge)
	return edge
}

func generateEdgeID(e *internal.Edge) string {
	return fmt.Sprintf("%s_%d_%d", e.Line.ID, e.From.EvaNumber, e.Actual.Departure.Unix())
}

func (c *container) insertStationEdge(last *internal.Edge, this *internal.Edge) *EdgePath {
	if last.To != this.From {
		log.Print("Tried to create stationEdge for line segments of different stations ", last.To.EvaNumber, this.From.EvaNumber)
		return nil
	}
	edge := &EdgePath{
		ID:   generateStationEdgeID(last, this),
		From: Coord{SpaceAxis: c.Stations[this.From], TimeAxis: last.Actual.Arrival},
		To:   Coord{SpaceAxis: c.Stations[this.From], TimeAxis: this.Actual.Departure},
		Edge: internal.Edge{
			Redundant: last.Redundant || this.Redundant,
		},
	}
	c.Edges[edge.ID] = edge
	c.SortedEdges = append(c.SortedEdges, edge)
	return edge
}

func generateStationEdgeID(last *internal.Edge, this *internal.Edge) string {
	if last == nil {
		return "undefined"
	}
	return fmt.Sprintf("%s_%d_%d_station", this.Line.ID, this.From.EvaNumber, last.Actual.Arrival.Unix())
}

func (c *container) gravitate() {
	c.TimeAxisSize = timeAxisSize
	c.SpaceAxisSize = spaceAxisSize

	c.layoutStations()
	c.indicateTimes()
}

func (c *container) layoutStations() {
	var stationsSlice []*StationLabel
	for _, s := range c.Stations {
		stationsSlice = append(stationsSlice, s)
	}
	sort.Slice(stationsSlice, func(i, j int) bool {
		return stationsSlice[i].Rank < stationsSlice[j].Rank
	})
	x := -1
	y := 0
	var lastGroup *int
	for _, s := range stationsSlice {
		log.Print(s.Name, s.GroupNumber, s.Rank)
		if s.GroupNumber != nil {
			log.Print(*s.GroupNumber)
		}
		if s.GroupNumber == nil || lastGroup == nil || *lastGroup != *s.GroupNumber {
			x++
			y = 0
			lastGroup = s.GroupNumber
		}
		s.SpaceAxis = x
		s.SpaceAxisHeap = y
		y++
	}
	c.maxSpace = x
}

func (c *container) indicateTimes() {
	delta := c.maxTime.Unix() - c.minTime.Unix()
	c.timeAxisDistance = float32(delta)
	//duration, _ := time.ParseDuration(fmt.Sprintf("%ds", delta/maxTimeIndicators))
	//now := c.minTime
}

func (c *container) render(wr io.Writer) {
	t := template.Must(template.New("time-space.tmpl.svg").ParseFiles("./render/time-space.tmpl.svg"))
	err := t.Execute(wr, c)
	if err != nil {
		panic(err)
	}
}

func (c *container) X(coord Coord) int {
	return int(float32(coord.SpaceAxis.SpaceAxis)/float32(c.maxSpace)*float32(c.SpaceAxisSize-50) + 50.0)
}

func (c *container) Y(coord Coord) int {
	if coord.TimeAxis.IsZero() {
		return 50 + coord.SpaceAxis.SpaceAxisHeap*20
	}
	delta := float32(coord.TimeAxis.Unix() - c.minTime.Unix())
	return int(delta/c.timeAxisDistance*float32(c.TimeAxisSize-100) + 100.0)
}

func (p *EdgePath) Label() string {
	e := p.Edge
	if e.Line == nil {
		return ""
	}
	var label string
	if e.Line.Name != "" {
		label = e.Line.Name
	} else {
		label = e.Line.ID
	}
	if e.Line.Type == "Foot" {
		return "🚶 " + label
	}
	return e.Line.Type + " " + label
}

func (p *EdgePath) Type() string {
	e := p.Edge
	if e.Line == nil {
		return ""
	}
	return e.Line.Type
}

func (p *EdgePath) Departure() string {
	e := p.Edge
	if e.Line == nil {
		return ""
	}
	var label string
	label = simpleTime(e.Actual.Departure) + delay(e.Current.Departure, e.Planned.Departure)
	if e.Planned.DepartureTrack != "" {
		label += e.Planned.DepartureTrack
	}
	return label
}

func (p *EdgePath) Arrival() string {
	e := p.Edge
	if e.Line == nil {
		return ""
	}
	label := simpleTime(e.Actual.Arrival) + delay(e.Current.Arrival, e.Planned.Arrival)
	if e.Planned.ArrivalTrack != "" {
		label += e.Planned.ArrivalTrack
	}
	return label
}

func simpleTime(t time.Time) string {
	return fmt.Sprintf("%02d:%02d ", t.Hour(), t.Minute())
}

func delay(current time.Time, planned time.Time) string {
	if !current.IsZero() {
		return " (" + fmt.Sprintf("%+.0f", current.Sub(planned).Minutes()) + ")"
	}
	return ""
}
