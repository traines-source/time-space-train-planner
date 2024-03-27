package internal

import "time"

type Station struct {
	Name       string
	ID         string
	Departures []*Edge
	Arrivals   []*Edge
	Lat        float64
	Lon        float64
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
	DestinationArrival         Distribution
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

type Distribution struct {
	Histogram           []float32
	Start               time.Time
	Mean                time.Time
	FeasibleProbability float32
}
