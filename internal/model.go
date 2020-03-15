package internal

import "time"

type Station struct {
	Name       string
	EvaNumber  int
	Departures []*Edge
	Lat        float32
	Lon        float32
}

type Edge struct {
	Line         *Line
	From         *Station
	To           *Station
	Planned      StopInfo
	Current      StopInfo
	Actual       StopInfo
	Message      string
	ShortestPath *Edge
	Redundant    bool
}

type StopInfo struct {
	Departure      time.Time
	Arrival        time.Time
	DepartureTrack string
}

type Line struct {
	Type    string
	Name    string
	ID      int
	Message string
	Route   []*Edge
	Stops   map[*Station]*LineStop
}

type LineStop struct {
	Station *Station
	Planned StopInfo
	Current StopInfo
	Message string
}
