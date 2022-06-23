package render

import (
	"fmt"
	"html"
	"io"
	"log"
	"math"
	"sort"
	"text/template"
	"time"

	"traines.eu/time-space-train-planner/internal"
)

type container struct {
	Stations            map[*internal.Station]*StationLabel
	Edges               map[string]*EdgePath
	SortedEdges         []*EdgePath
	MinTime             time.Time
	MaxTime             time.Time
	maxSpace            int
	timeAxisDistance    float32
	TimeIndicators      []Coord
	TimeAxisSize        int
	SpaceAxisSize       int
	Query               string
	DefaultShortestPathID string
}

const (
	timeAxisSize  = 1500
	spaceAxisSize = 1500
)

func TimeSpace(stations map[int]*internal.Station, lines map[string]*internal.Line, wr io.Writer, query string) {
	c := &container{Stations: map[*internal.Station]*StationLabel{}, Edges: map[string]*EdgePath{}, Query: html.EscapeString(query)}
	c.setupStations(stations)
	c.setupEdges(lines)
	c.setupPreviousAndNext(stations)
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
					c.setShortestPathFor(originEdgePath, e, lastEdge, e)
					lastEdge = e
				}
				for e := origin; e != nil; e = e.ReverseShortestPath {
					c.setShortestPathFor(originEdgePath, e, e, lastEdge)
					lastEdge = e
				}
			}
		}
	}
}

func (c *container) setupPreviousAndNext(stations map[int]*internal.Station) {
	var stationsSlice []*internal.Station
	for _, s := range stations {
		stationsSlice = append(stationsSlice, s)
	}
	sort.Slice(stationsSlice, func(i, j int) bool {
		return stationsSlice[i].Rank < stationsSlice[j].Rank
	})
	c.preselectShortestPath(stationsSlice[len(stationsSlice)-1])
	var arrivals []*internal.Edge
	var departures []*internal.Edge
	var lastGroup *int
	for _, s := range stationsSlice {
		if s.GroupNumber == nil || lastGroup == nil || *lastGroup != *s.GroupNumber {
			c.flushStationGroup(stations, departures, arrivals)
			arrivals = []*internal.Edge{}
			departures = []*internal.Edge{}
			lastGroup = s.GroupNumber
		}
		arrivals = append(arrivals, s.Arrivals...)
		departures = append(departures, s.Departures...)
	}
	c.flushStationGroup(stations, departures, arrivals)
}

func (c *container) flushStationGroup(stations map[int]*internal.Station, departures []*internal.Edge, arrivals []*internal.Edge) {
	sort.SliceStable(arrivals, func(i, j int) bool {
		return arrivals[i].Actual.Arrival.Before(arrivals[j].Actual.Arrival)
	})
	sort.SliceStable(departures, func(i, j int) bool {
		return departures[i].Actual.Departure.Before(departures[j].Actual.Departure)
	})
	lastProperArrival := ""
	lastProperDeparture := ""
	var nextArrivalToFill *string = nil
	var nextDepartureToFill *string = nil
	for i := 0; i < len(arrivals); i++ {
		if e, ok := c.Edges[generateEdgeID(arrivals[i])]; ok {
			if isEdgeInsideGroup(e) {
				continue
			}
			if nextArrivalToFill != nil {
				*nextArrivalToFill = "data-na=\"" + generateEdgeID(arrivals[i]) + "\""
				nextArrivalToFill = nil
			}
			e.PreviousArrival = "data-pa=\"" + lastProperArrival + "\""
			if i+1 < len(arrivals) && e.To.SpaceAxis.Station.GroupNumber != nil && stations[*e.To.SpaceAxis.Station.GroupNumber].Rank+1 == len(stations) {
				nextArrivalToFill = &e.NextArrival
			}
			lastProperArrival = generateEdgeID(arrivals[i])
		}
	}
	for i := 0; i < len(departures); i++ {
		if e, ok := c.Edges[generateEdgeID(departures[i])]; ok {
			if isEdgeInsideGroup(e) {
				continue
			}
			if nextDepartureToFill != nil {
				*nextDepartureToFill = "data-nd=\"" + generateEdgeID(departures[i]) + "\""
				nextDepartureToFill = nil
			}
			if e.From.SpaceAxis.Station.GroupNumber != nil && stations[*e.From.SpaceAxis.Station.GroupNumber].Rank == 0 {
				e.PreviousDeparture = "data-pd=\"" + lastProperDeparture + "\""
			}
			nextDepartureToFill = &e.NextDeparture
			lastProperDeparture = generateEdgeID(departures[i])
		}
	}
}

func (c *container) preselectShortestPath(destination *internal.Station) {
	for _, s := range destination.Arrivals {
		if s.ReverseShortestPath != nil {
			if e, ok := c.Edges[generateEdgeID(s)]; ok {
				c.DefaultShortestPathID = e.ID
			}
			break
		}
	}
}

func isEdgeInsideGroup(e *EdgePath) bool {
	return e.From.SpaceAxis.GroupNumber != nil && e.To.SpaceAxis.GroupNumber != nil && *e.From.SpaceAxis.GroupNumber == *e.To.SpaceAxis.GroupNumber
}

func (c *container) setShortestPathFor(originEdgePath *EdgePath, e *internal.Edge, start *internal.Edge, end *internal.Edge) {
	if edgePath, ok := c.Edges[generateEdgeID(e)]; ok {
		edgePath.ShortestPathFor = append(edgePath.ShortestPathFor, originEdgePath)
	}
	if edgePath, ok := c.Edges[generateStationEdgeID(start, end)]; ok {
		edgePath.ShortestPathFor = append(edgePath.ShortestPathFor, originEdgePath)
	}
}

func (c *container) stretchTimeAxis(min time.Time, max time.Time) {
	if min.Before(c.MinTime) || c.MinTime.IsZero() {
		c.MinTime = min
	}
	if max.After(c.MaxTime) || c.MaxTime.IsZero() {
		c.MaxTime = max
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
	return fmt.Sprintf("%p", e)
}

func (c *container) insertStationEdge(last *internal.Edge, this *internal.Edge) *EdgePath {
	if last.To != this.From {
		log.Print("Tried to create stationEdge for line segments of different stations ", last.To.EvaNumber, this.From.EvaNumber)
		return nil
	}
	if this.Actual.Departure.Sub(last.Actual.Arrival).Minutes() > 30 {
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
	return fmt.Sprintf("%p_%p_station", last, this)
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
		if s.GroupNumber == nil || lastGroup == nil || *lastGroup != *s.GroupNumber {
			x++
			y = 0
			lastGroup = s.GroupNumber
		}
		s.SpaceAxis = x
		if s.GroupNumber == nil || s.GroupNumber != nil && *s.GroupNumber == s.EvaNumber {
			s.SpaceAxisHeap = 0
		} else {
			y++
			s.SpaceAxisHeap = y
		}
	}
	c.maxSpace = x
}

func (c *container) indicateTimes() {
	delta := c.MaxTime.Unix() - c.MinTime.Unix()
	c.timeAxisDistance = float32(delta)
	t := time.Date(c.MinTime.Year(), c.MinTime.Month(), c.MinTime.Day(), c.MinTime.Hour()+1, 0, 0, 0, c.MinTime.Location())
	for ; t.Before(c.MaxTime); t = t.Add(time.Hour) {
		c.TimeIndicators = append(c.TimeIndicators, Coord{TimeAxis: t})
	}
}

func (c *container) render(wr io.Writer) {
	t := template.Must(template.New("time-space.tmpl.svg").ParseFiles("./render/time-space.tmpl.svg"))
	err := t.Execute(wr, c)
	if err != nil {
		panic(err)
	}
}

func (c *container) X(coord Coord) int {
	if coord.SpaceAxis == nil {
		return 0
	}
	return int(float32(coord.SpaceAxis.SpaceAxis)/float32(c.maxSpace)*float32(c.SpaceAxisSize-100) + 50.0)
}

func (c *container) Y(coord Coord) int {
	if coord.TimeAxis.IsZero() {
		return 50 + coord.SpaceAxis.SpaceAxisHeap*20
	}
	delta := float32(coord.TimeAxis.Unix() - c.MinTime.Unix())
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
	if e.Message != "" {
		label += " (" + e.Message[0:int(math.Min(float64(len(e.Message)), 30))] + "...)"
	}
	if e.Line.Type == "Foot" {
		return "🚶 " + label
	}
	return label
}

func (p *EdgePath) Type() string {
	e := p.Edge
	if e.Line == nil {
		return ""
	}
	return e.Line.Type
}

func (p *EdgePath) Departure() string {
	return p.time(func(stop internal.StopInfo) time.Time { return stop.Departure }, func(stop internal.StopInfo) string { return stop.DepartureTrack })
}

func (p *EdgePath) Arrival() string {
	return p.time(func(stop internal.StopInfo) time.Time { return stop.Arrival }, func(stop internal.StopInfo) string { return stop.ArrivalTrack })
}

func (p *EdgePath) LiveDataDeparture() string {
	return p.liveDataClass(func(stop internal.StopInfo) time.Time { return stop.Departure })
}

func (p *EdgePath) LiveDataArrival() string {
	return p.liveDataClass(func(stop internal.StopInfo) time.Time { return stop.Arrival })
}

func (p *EdgePath) time(timeResolver func(internal.StopInfo) time.Time, trackResolver func(internal.StopInfo) string) string {
	e := p.Edge
	if e.Line == nil {
		return ""
	}
	label := simpleTime(timeResolver(e.Actual)) + delay(timeResolver(e.Current), timeResolver(e.Planned))
	if trackResolver(e.Planned) != "" {
		label += "Pl." + trackResolver(e.Planned)
	}
	return label
}

func (p *EdgePath) liveDataClass(timeResolver func(internal.StopInfo) time.Time) string {
	e := p.Edge
	if e.Line == nil {
		return ""
	}
	current := timeResolver(e.Current)
	if current.IsZero() {
		return ""
	}
	if current.Sub(timeResolver(e.Planned)).Minutes() > 5 {
		return "live-red"
	}
	return "live-green"
}

func (c *container) Minutes(time time.Time) string {
	return fmt.Sprintf("%.0f", time.Sub(c.MinTime).Minutes())
}

func simpleTime(t time.Time) string {
	return fmt.Sprintf("%02d:%02d ", t.Hour(), t.Minute())
}

func (c *container) SimpleTime(t time.Time) string {
	return simpleTime(t)
}

func delay(current time.Time, planned time.Time) string {
	if !current.IsZero() {
		return " (" + fmt.Sprintf("%+.0f", current.Sub(planned).Minutes()) + ") "
	}
	return ""
}
