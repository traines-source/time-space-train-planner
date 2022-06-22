package internal

import "time"

type Station struct {
	Name        string
	EvaNumber   int
	Departures  []*Edge
	Arrivals    []*Edge
	Lat         float32
	Lon         float32
	Rank        int
	GroupNumber *int
}

type Edge struct {
	Line                 *Line
	From                 *Station
	To                   *Station
	Planned              StopInfo
	Current              StopInfo
	Actual               StopInfo
	Message              string
	ShortestPath         *Edge
	ReverseShortestPath  *Edge
	ProviderShortestPath bool
	ShortestPathFor      map[*Edge]struct{}
	Redundant            bool
	Discarded            bool
}

type StopInfo struct {
	Departure      time.Time
	Arrival        time.Time
	DepartureTrack string
	ArrivalTrack   string
}

type Line struct {
	Type    string
	Name    string
	ID      string
	Message string
	Route   []*Edge
	Stops   []*LineStop
}

type LineStop struct {
	Station *Station
	Planned StopInfo
	Current StopInfo
	Message string
}
