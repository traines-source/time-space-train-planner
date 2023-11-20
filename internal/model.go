package internal

import "time"

type Station struct {
	Name       string
	ID         string
	Departures []*Edge
	Arrivals   []*Edge
	Lat        float32
	Lon        float32
	Rank       int
	GroupID    *string
}

type Edge struct {
	Line                       *Line
	From                       *Station
	To                         *Station
	Planned                    StopInfo
	Current                    StopInfo
	Actual                     StopInfo
	Message                    string
	ShortestPath               *Edge
	ReverseShortestPath        *Edge
	ProviderShortestPath       bool
	ShortestPathFor            map[*Edge]struct{}
	EarliestDestinationArrival time.Time
	LatestOriginDeparture      time.Time
	Redundant                  bool
	Discarded                  bool
	Cancelled                  bool
}

type StopInfo struct {
	Departure      time.Time
	Arrival        time.Time
	DepartureTrack string
	ArrivalTrack   string
}

type Line struct {
	Type      string
	Name      string
	ID        string
	Message   string
	Direction string
	Route     []*Edge
	Stops     []*LineStop
}

type LineStop struct {
	Station   *Station
	Planned   StopInfo
	Current   StopInfo
	Message   string
	Cancelled bool
}
