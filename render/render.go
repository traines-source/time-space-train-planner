package render

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"text/template"
	"time"

	"traines.eu/time-space-train-planner/internal"
)

type container struct {
	From                  StationLabel
	To                    StationLabel
	Vias                  []StationLabel
	Stations              map[string]*StationLabel
	Edges                 map[string]*EdgePath
	SortedEdges           []string
	MinTime               time.Time
	MaxTime               time.Time
	MaxSpace              int
	TimeAxisDistance      float32
	TimeIndicators        []Coord
	TimeAxisSize          int
	SpaceAxisSize         int
	Query                 string
	DefaultShortestPathID string
	LegalLink             string
}

const (
	timeAxisSize  = 1500
	spaceAxisSize = 1500
)

func TimeSpace(stations map[int]*internal.Station, lines map[string]*internal.Line, wr io.Writer, query string) {
	c := &container{
		Stations:  map[string]*StationLabel{},
		Edges:     map[string]*EdgePath{},
		Query:     html.EscapeString(query),
		LegalLink: os.Getenv("TSTP_LEGAL"),
	}
	c.setupStations(stations)
	c.setupEdges(lines)
	c.setupShortestPathFors(lines)
	c.setupPreviousAndNext(stations)
	c.gravitate()
	c.render(wr)
}

func TimeSpaceApi(stations map[int]*internal.Station, lines map[string]*internal.Line, wr io.Writer, query string) {
	c := &container{
		Stations: map[string]*StationLabel{},
		Edges:    map[string]*EdgePath{},
	}
	c.setupStations(stations)
	c.setupEdges(lines)
	c.setupPreviousAndNext(stations)
	c.gravitate()
	json.NewEncoder(wr).Encode(c)
}

func (c *container) setupStations(stations map[int]*internal.Station) {
	var from int
	var to int
	for _, s := range stations {
		if len(s.Arrivals) > 0 || len(s.Departures) > 0 {
			station := &StationLabel{
				ID:   strconv.Itoa(s.EvaNumber),
				Name: s.Name,
				Rank: s.Rank,
			}
			if s.GroupNumber != nil {
				g := strconv.Itoa(*s.GroupNumber)
				station.GroupID = &g
			}
			if s.Rank == 0 {
				from = s.EvaNumber
				c.From = *station
			}
			if s.Rank+1 == len(stations) {
				to = s.EvaNumber
				c.To = *station
			}
			station.Coord.SpaceAxis = station.ID
			c.Stations[station.ID] = station
		}
	}
	c.Vias = makeVias(stations, from, to)
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
		if c.Edges[c.SortedEdges[i]].Redundant == c.Edges[c.SortedEdges[j]].Redundant && c.Edges[c.SortedEdges[i]].Line != nil {
			return c.Edges[c.SortedEdges[i]].Line.Type == "Foot"
		}
		return c.Edges[c.SortedEdges[i]].Redundant
	})
}

func (c *container) setupShortestPathFors(lines map[string]*internal.Line) {
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
	c.preselectShortestPath(stationsSlice[0], stationsSlice[len(stationsSlice)-1])
	var arrivals []*internal.Edge
	var departures []*internal.Edge
	var lastGroup *int
	for _, s := range stationsSlice {
		if s.GroupNumber == nil || lastGroup == nil || *lastGroup != *s.GroupNumber {
			c.flushStationGroup(departures, arrivals)
			arrivals = []*internal.Edge{}
			departures = []*internal.Edge{}
			lastGroup = s.GroupNumber
		}
		arrivals = append(arrivals, s.Arrivals...)
		departures = append(departures, s.Departures...)
	}
	c.flushStationGroup(departures, arrivals)
}

func (c *container) flushStationGroup(departures []*internal.Edge, arrivals []*internal.Edge) {
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
			if c.isEdgeInsideGroup(e) {
				continue
			}
			if e.Discarded {
				continue
			}
			if nextArrivalToFill != nil {
				*nextArrivalToFill = generateEdgeID(arrivals[i])
				nextArrivalToFill = nil
			}
			e.PreviousArrival = lastProperArrival
			if i+1 < len(arrivals) && c.Stations[e.To.SpaceAxis].GroupID != nil {
				if val, ok := c.Stations[*c.Stations[e.To.SpaceAxis].GroupID]; ok && val.Rank+1 == len(c.Stations) {
					nextArrivalToFill = &e.NextArrival
				}
			}
			lastProperArrival = generateEdgeID(arrivals[i])
		}
	}
	for i := 0; i < len(departures); i++ {
		if e, ok := c.Edges[generateEdgeID(departures[i])]; ok {
			if c.isEdgeInsideGroup(e) {
				continue
			}
			if e.Discarded {
				continue
			}
			if nextDepartureToFill != nil {
				*nextDepartureToFill = generateEdgeID(departures[i])
				nextDepartureToFill = nil
			}
			if c.Stations[e.From.SpaceAxis].GroupID != nil {
				if val, ok := c.Stations[*c.Stations[e.From.SpaceAxis].GroupID]; ok && val.Rank == 0 {
					e.PreviousDeparture = lastProperDeparture
				}
			}
			nextDepartureToFill = &e.NextDeparture
			lastProperDeparture = generateEdgeID(departures[i])
		}
	}
}

func (c *container) preselectShortestPath(origin *internal.Station, destination *internal.Station) {
	for _, s := range destination.Arrivals {
		if s.ReverseShortestPath != nil || s.From.EvaNumber == origin.EvaNumber {
			start := s
			for start.ReverseShortestPath != nil {
				start = s.ReverseShortestPath
			}
			if e, ok := c.Edges[generateEdgeID(start)]; ok {
				c.DefaultShortestPathID = e.ID
			}
			break
		}
	}
}

func (c *container) isEdgeInsideGroup(e *EdgePath) bool {
	return c.Stations[e.From.SpaceAxis].GroupID != nil && c.Stations[e.To.SpaceAxis].GroupID != nil && *c.Stations[e.From.SpaceAxis].GroupID == *c.Stations[e.To.SpaceAxis].GroupID
}

func (c *container) setShortestPathFor(originEdgePath *EdgePath, e *internal.Edge, start *internal.Edge, end *internal.Edge) {
	if edgePath, ok := c.Edges[generateEdgeID(e)]; ok {
		edgePath.ShortestPathFor = append(edgePath.ShortestPathFor, originEdgePath.ID)
	}
	if edgePath, ok := c.Edges[generateStationEdgeID(start, end)]; ok {
		edgePath.ShortestPathFor = append(edgePath.ShortestPathFor, originEdgePath.ID)
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
		ID:                   generateEdgeID(e),
		From:                 Coord{SpaceAxis: strconv.Itoa(e.From.EvaNumber), TimeAxis: e.Actual.Departure},
		To:                   Coord{SpaceAxis: strconv.Itoa(e.To.EvaNumber), TimeAxis: e.Actual.Arrival},
		Redundant:            e.Redundant,
		Discarded:            e.Discarded,
		Cancelled:            e.Cancelled,
		Message:              e.Message,
		Planned:              e.Planned,
		Current:              e.Current,
		Actual:               e.Actual,
		ShortestPath:         []ShortestPathAlternative{},
		ReverseShortestPath:  []ShortestPathAlternative{},
		ProviderShortestPath: e.ProviderShortestPath,
		ShortestPathFor:      []string{},
	}
	if e.Line != nil {
		edge.Line = &LineLabel{
			Name:      e.Line.Name,
			ID:        e.Line.ID,
			Type:      e.Line.Type,
			Direction: e.Line.Direction,
		}
	}
	if e.ShortestPath != nil {
		edge.ShortestPath = append(edge.ShortestPath, ShortestPathAlternative{EdgeID: generateEdgeID(e.ShortestPath)})
	}
	if e.ReverseShortestPath != nil {
		edge.ReverseShortestPath = append(edge.ReverseShortestPath, ShortestPathAlternative{EdgeID: generateEdgeID(e.ReverseShortestPath)})
	}
	c.Edges[edge.ID] = edge
	c.SortedEdges = append(c.SortedEdges, edge.ID)
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
		ID:                  generateStationEdgeID(last, this),
		ShortestPathFor:     []string{},
		ShortestPath:        []ShortestPathAlternative{},
		ReverseShortestPath: []ShortestPathAlternative{},
		From:                Coord{SpaceAxis: strconv.Itoa(this.From.EvaNumber), TimeAxis: last.Actual.Arrival},
		To:                  Coord{SpaceAxis: strconv.Itoa(this.From.EvaNumber), TimeAxis: this.Actual.Departure},
		Redundant:           last.Redundant || this.Redundant,
	}
	c.Edges[edge.ID] = edge
	c.SortedEdges = append(c.SortedEdges, edge.ID)
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
	var lastGroup *string
	for _, s := range stationsSlice {
		if s.GroupID == nil || lastGroup == nil || *lastGroup != *s.GroupID {
			x++
			y = 0
			lastGroup = s.GroupID
		}
		s.SpaceAxis = x
		if s.GroupID == nil || s.GroupID != nil && *s.GroupID == s.ID {
			s.SpaceAxisHeap = 0
		} else {
			y++
			s.SpaceAxisHeap = y
		}
	}
	c.MaxSpace = x
}

func (c *container) indicateTimes() {
	delta := c.MaxTime.Unix() - c.MinTime.Unix()
	c.TimeAxisDistance = float32(delta)
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
	if coord.SpaceAxis == "" {
		return 0
	}
	return int(float32(c.Stations[coord.SpaceAxis].SpaceAxis)/float32(c.MaxSpace)*float32(c.SpaceAxisSize-100) + 50.0)
}

func (c *container) Y(coord Coord) int {
	if coord.TimeAxis.IsZero() {
		return 50 + c.Stations[coord.SpaceAxis].SpaceAxisHeap*20
	}
	delta := float32(coord.TimeAxis.Unix() - c.MinTime.Unix())
	return int(delta/c.TimeAxisDistance*float32(c.TimeAxisSize-100) + 100.0)
}

func (e *EdgePath) Label() string {
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
		label += " (" + substr(e.Message, 0, 30) + "...)"
	}
	if e.Line.Type == "Foot" {
		return "🚶 " + label
	}
	return label
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func (e *EdgePath) Type() string {
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

func (e *EdgePath) time(timeResolver func(internal.StopInfo) time.Time, trackResolver func(internal.StopInfo) string) string {
	if e.Line == nil {
		return ""
	}
	label := simpleTime(timeResolver(e.Actual)) + delay(timeResolver(e.Current), timeResolver(e.Planned))
	if trackResolver(e.Planned) != "" {
		label += "Pl." + trackResolver(e.Planned)
	}
	return label
}

func (e *EdgePath) liveDataClass(timeResolver func(internal.StopInfo) time.Time) string {
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
