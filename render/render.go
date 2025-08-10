package render

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"os"
	"sort"
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
	edgeIDByEdge          map[*internal.Edge]string
	edgeByEdgeID          map[string]*internal.Edge
}

const (
	timeAxisSize  = 1500
	spaceAxisSize = 1500
)

func TimeSpace(stations map[string]*internal.Station, lines map[string]*internal.Line, wr io.Writer, query string) {
	c := &container{
		Stations:     map[string]*StationLabel{},
		Edges:        map[string]*EdgePath{},
		edgeIDByEdge: map[*internal.Edge]string{},
		edgeByEdgeID: map[string]*internal.Edge{},
		Query:        html.EscapeString(query),
		LegalLink:    os.Getenv("TSTP_LEGAL"),
	}
	c.setupStations(stations)
	c.setupEdges(lines)
	c.setupPreviousAndNext(stations)
	c.gravitate()
	c.render(wr)
}

func TimeSpaceApi(stations map[string]*internal.Station, lines map[string]*internal.Line, wr io.Writer, query string) {
	c := &container{
		Stations:     map[string]*StationLabel{},
		Edges:        map[string]*EdgePath{},
		edgeIDByEdge: map[*internal.Edge]string{},
		edgeByEdgeID: map[string]*internal.Edge{},
	}
	c.setupStations(stations)
	c.setupEdges(lines)
	c.setupPreviousAndNext(stations)
	c.gravitate()
	log.Print("Done.")
	json.NewEncoder(wr).Encode(c)
}

func (c *container) setupStations(stations map[string]*internal.Station) {
	var from string
	var to string
	for _, s := range stations {
		if len(s.Arrivals) > 0 || len(s.Departures) > 0 {
			station := &StationLabel{
				ID:             s.ID,
				Name:           s.Name,
				Rank:           s.Rank,
				BestDepartures: []string{},
				Lon:            s.Lon,
				Lat:            s.Lat,
			}
			station.GroupID = s.GroupID
			if s.Rank == 0 {
				from = s.ID
				c.From = *station
			}
			if s.Rank+1 == len(stations) {
				to = s.ID
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
			if e.Discarded {
				continue
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
	for _, s := range c.Stations {
		sort.SliceStable(s.BestDepartures, func(i, j int) bool {
			return c.Edges[s.BestDepartures[i]].Planned.Departure.Before(c.Edges[s.BestDepartures[j]].Planned.Departure)
		})
		sort.SliceStable(s.BestDepartures, func(i, j int) bool {
			return c.Edges[s.BestDepartures[i]].EarliestDestinationArrival.Before(c.Edges[s.BestDepartures[j]].EarliestDestinationArrival)
		})
		sort.SliceStable(s.BestDepartures, func(i, j int) bool {
			a := c.Edges[s.BestDepartures[i]].DestinationArrival.Mean
			b := c.Edges[s.BestDepartures[j]].DestinationArrival.Mean
			return !a.IsZero() && a.Before(b)
		})
	}
}

func (c *container) setupPreviousAndNext(stations map[string]*internal.Station) {
	var stationsSlice []*internal.Station
	for _, s := range stations {
		stationsSlice = append(stationsSlice, s)
	}
	sort.Slice(stationsSlice, func(i, j int) bool {
		return stationsSlice[i].Rank < stationsSlice[j].Rank
	})
	var arrivals []*internal.Edge
	var departures []*internal.Edge
	var lastGroup *string
	for _, s := range stationsSlice {
		if s.GroupID == nil || lastGroup == nil || *lastGroup != *s.GroupID {
			c.flushStationGroup(departures, arrivals)
			arrivals = []*internal.Edge{}
			departures = []*internal.Edge{}
			lastGroup = s.GroupID
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
		if e, ok := c.Edges[c.generateEdgeID(arrivals[i])]; ok {
			if c.isEdgeInsideGroup(e) {
				continue
			}
			if e.Discarded {
				continue
			}
			if nextArrivalToFill != nil {
				*nextArrivalToFill = c.generateEdgeID(arrivals[i])
				nextArrivalToFill = nil
			}
			e.PreviousArrival = lastProperArrival
			if i+1 < len(arrivals) && c.Stations[e.To.SpaceAxis].GroupID != nil {
				if val, ok := c.Stations[*c.Stations[e.To.SpaceAxis].GroupID]; ok && val.Rank+1 == len(c.Stations) {
					nextArrivalToFill = &e.NextArrival
				}
			}
			lastProperArrival = c.generateEdgeID(arrivals[i])
		} else if !arrivals[i].Discarded {
			log.Print("Referenced non-existing edge. (arrows1)")
		}
	}
	for i := 0; i < len(departures); i++ {
		if e, ok := c.Edges[c.generateEdgeID(departures[i])]; ok {
			if c.isEdgeInsideGroup(e) {
				continue
			}
			if e.Discarded {
				continue
			}
			if nextDepartureToFill != nil {
				*nextDepartureToFill = c.generateEdgeID(departures[i])
				nextDepartureToFill = nil
			}
			if c.Stations[e.From.SpaceAxis].GroupID != nil {
				if val, ok := c.Stations[*c.Stations[e.From.SpaceAxis].GroupID]; ok && val.Rank == 0 {
					e.PreviousDeparture = lastProperDeparture
				}
			}
			nextDepartureToFill = &e.NextDeparture
			lastProperDeparture = c.generateEdgeID(departures[i])
		} else if !departures[i].Discarded {
			log.Print("Referenced non-existing edge. (arrows2)")
		}
	}
}

func (c *container) isEdgeInsideGroup(e *EdgePath) bool {
	return c.Stations[e.From.SpaceAxis].GroupID != nil && c.Stations[e.To.SpaceAxis].GroupID != nil && *c.Stations[e.From.SpaceAxis].GroupID == *c.Stations[e.To.SpaceAxis].GroupID
}

func (c *container) insertEdge(e *internal.Edge) *EdgePath {
	edgeID := c.generateEdgeID(e)
	edge := &EdgePath{
		ID:                         edgeID,
		From:                       Coord{SpaceAxis: e.From.ID, TimeAxis: e.Actual.Departure},
		To:                         Coord{SpaceAxis: e.To.ID, TimeAxis: e.Actual.Arrival},
		Redundant:                  e.Redundant,
		Discarded:                  e.Discarded,
		Cancelled:                  e.Cancelled,
		Message:                    e.Message,
		Planned:                    e.Planned,
		Current:                    e.Current,
		Actual:                     e.Actual,
		ShortestPath:               []ShortestPathAlternative{},
		ReverseShortestPath:        []ShortestPathAlternative{},
		ProviderShortestPath:       e.ProviderShortestPath,
		ShortestPathFor:            []string{},
		EarliestDestinationArrival: e.EarliestDestinationArrival,
		DestinationArrival:         e.DestinationArrival,
	}
	if e.Line != nil {
		edge.Line = &LineLabel{
			ID:        e.Line.ID,
			Name:      e.Line.Name,
			Type:      e.Line.Type,
			Direction: e.Line.Direction,
		}
	}
	if e.ShortestPath != nil {
		edge.ShortestPath = append(edge.ShortestPath, ShortestPathAlternative{EdgeID: c.generateEdgeID(e.ShortestPath)})
	}
	if e.ReverseShortestPath != nil {
		edge.ReverseShortestPath = append(edge.ReverseShortestPath, ShortestPathAlternative{EdgeID: c.generateEdgeID(e.ReverseShortestPath)})
	}
	if !e.Discarded {
		if station, ok := c.Stations[e.From.ID]; ok && !edge.EarliestDestinationArrival.IsZero() {
			station.BestDepartures = append(station.BestDepartures, edgeID)
		}
	}
	c.Edges[edge.ID] = edge
	c.SortedEdges = append(c.SortedEdges, edge.ID)
	return edge
}

func (c *container) generateEdgeID(e *internal.Edge) string {
	if edgeID, ok := c.edgeIDByEdge[e]; ok {
		return edgeID
	}
	hash := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s_%s_%s_%d", e.Line.ID, e.From.ID, e.To.ID, e.Planned.Departure.Unix()))))[:7]
	for {
		if _, ok := c.edgeByEdgeID[hash]; !ok {
			c.edgeIDByEdge[e] = hash
			c.edgeByEdgeID[hash] = e
			return hash
		}
		log.Print("Prevented id collision for line", e.Line.ID, hash)
		hash += "C"
	}
}

func (c *container) generateStationEdgeID(last *internal.Edge, this *internal.Edge) string {
	if last == nil || this == nil {
		return "undefined"
	}
	return fmt.Sprintf("%s_%s_station", c.generateEdgeID(last), c.generateEdgeID(this))
}

func (c *container) insertStationEdge(last *internal.Edge, this *internal.Edge) *EdgePath {
	if this.Line.Type == "Foot" {
		return nil
	}
	if last.To != this.From {
		log.Print("Tried to create stationEdge for line segments of different stations ", last.To.ID, this.From.ID, last.Message)
		return nil
	}
	if this.Actual.Departure.Sub(last.Actual.Arrival).Minutes() > 30 {
		return nil
	}
	edge := &EdgePath{
		ID:                  c.generateStationEdgeID(last, this),
		ShortestPathFor:     []string{},
		ShortestPath:        []ShortestPathAlternative{},
		ReverseShortestPath: []ShortestPathAlternative{},
		From:                Coord{SpaceAxis: this.From.ID, TimeAxis: last.Actual.Arrival},
		To:                  Coord{SpaceAxis: this.From.ID, TimeAxis: this.Actual.Departure},
		Redundant:           last.Redundant || this.Redundant,
	}
	c.Edges[edge.ID] = edge
	c.SortedEdges = append(c.SortedEdges, edge.ID)
	return edge
}

func (c *container) gravitate() {
	c.TimeAxisSize = timeAxisSize
	c.SpaceAxisSize = spaceAxisSize

	c.layoutStations()
	c.indicateTimes()
}

func (c *container) stretchTimeAxis(min time.Time, max time.Time) {
	if min.Before(c.MinTime) || c.MinTime.IsZero() {
		c.MinTime = min
	}
	if max.After(c.MaxTime) || c.MaxTime.IsZero() {
		c.MaxTime = max
	}
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
	for i, s := range stationsSlice {
		if s.GroupID == nil || lastGroup == nil || *lastGroup != *s.GroupID {
			x++
			y = 0
			lastGroup = s.GroupID
		}
		s.SpaceAxis = x
		if s.GroupID == nil || s.GroupID != nil && *s.GroupID == s.ID {
			s.SpaceAxisHeap = 0
		} else {
			if i == 0 || stationsSlice[i-1].Name != s.Name {
				y++
			}
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
