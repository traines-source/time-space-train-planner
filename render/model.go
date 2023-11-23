package render

import (
	"sort"
	"time"

	"traines.eu/time-space-train-planner/internal"
)

type StationLabel struct {
	ID             string
	Name           string
	Coord          Coord
	SpaceAxis      int
	SpaceAxisHeap  int
	GroupID        *string
	Rank           int
	BestDepartures []string
}

type EdgePath struct {
	ID                         string
	From                       Coord
	To                         Coord
	Redundant                  bool
	Discarded                  bool
	Cancelled                  bool
	Line                       *LineLabel
	Message                    string
	Planned                    internal.StopInfo
	Current                    internal.StopInfo
	Actual                     internal.StopInfo
	ShortestPath               []ShortestPathAlternative
	ReverseShortestPath        []ShortestPathAlternative
	ProviderShortestPath       bool
	ShortestPathFor            []string
	EarliestDestinationArrival time.Time
	DestinationArrival         internal.Distribution
	PreviousDeparture          string
	NextDeparture              string
	PreviousArrival            string
	NextArrival                string
}

type LineLabel struct {
	ID        string
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

func makeVias(stations map[string]*internal.Station, from string, to string) []StationLabel {
	list := []StationLabel{}
	for _, s := range stations {
		if s.ID == from || s.ID == to || s.GroupID != nil && *s.GroupID != s.ID {
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
	return StationLabel{ID: s.ID, Name: s.Name, Rank: s.Rank}
}
