package render

import (
	"time"

	"traines.eu/time-space-train-planner/internal"
)

type StationLabel struct {
	internal.Station
	Coord     Coord
	SpaceAxis int
}

type EdgePath struct {
	internal.Edge
	ID              string
	ShortestPathFor []*EdgePath
	From            Coord
	To              Coord
	Label           string
}

type Coord struct {
	TimeAxis  time.Time
	SpaceAxis *StationLabel
}
