package render

import (
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
	ShortestPathFor      []string
	From                 Coord
	To                   Coord
	Redundant            bool
	Discarded            bool
	Line                 *LineLabel
	Message              string
	Planned              internal.StopInfo
	Current              internal.StopInfo
	Actual               internal.StopInfo
	ProviderShortestPath bool
	PreviousDeparture    string
	NextDeparture        string
	PreviousArrival      string
	NextArrival          string
}

type LineLabel struct {
	Name string
	ID   string
	Type string
}

type Coord struct {
	TimeAxis  time.Time
	SpaceAxis string
}
