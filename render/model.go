package render

import (
	"sort"
	"strconv"
	"time"

	"traines.eu/time-space-train-planner/internal"
)

type StationLabel struct {
	ID            string
	Name          string
	Coord         Coord
	SpaceAxis     int
	SpaceAxisHeap int
	GroupID       *string
	Rank          int
}

type EdgePath struct {
	ID                   string
	From                 Coord
	To                   Coord
	Redundant            bool
	Discarded            bool
	Cancelled            bool
	Line                 *LineLabel
	Message              string
	Planned              internal.StopInfo
	Current              internal.StopInfo
	Actual               internal.StopInfo
	ShortestPath         []ShortestPathAlternative
	ReverseShortestPath  []ShortestPathAlternative
	ProviderShortestPath bool
	ShortestPathFor      []string
	PreviousDeparture    string
	NextDeparture        string
	PreviousArrival      string
	NextArrival          string
}

type LineLabel struct {
	Name      string
	Type      string
	Direction string
}

type ShortestPathAlternative struct {
	EdgeID string
}

type Coord struct {
	TimeAxis  time.Time
	SpaceAxis string
}

func makeVias(stations map[int]*internal.Station, from int, to int) []StationLabel {
	list := []StationLabel{}
	for _, s := range stations {
		if s.EvaNumber == from || s.EvaNumber == to || s.GroupNumber != nil && *s.GroupNumber != s.EvaNumber {
			continue
		}
		list = append(list, makeStationLabel(s))
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Rank < list[j].Rank
	})
	return list
}

func makeStationLabel(s *internal.Station) StationLabel {
	if s == nil {
		return StationLabel{}
	}
	return StationLabel{ID: strconv.Itoa(s.EvaNumber), Name: s.Name, Rank: s.Rank}
}
