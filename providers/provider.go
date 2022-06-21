package providers

import "time"

type Provider interface {
	Fetch(c Consumer)
	Enrich(c Consumer)
}

type Consumer interface {
	RequestStationDataBetween(station *ProviderStation) (from time.Time, to time.Time)
	Stations() []ProviderStation
	StationByName(name string) (ProviderStation, error)
	StationByEva(evaNumber int) (ProviderStation, error)
	UpsertStation(station ProviderStation)
	UpsertLine(line ProviderLine)
	UpsertLineStop(lineStop ProviderLineStop)
	UpsertLineEdge(lineEdge ProviderLineEdge)
	SetExpectedTravelDuration(duration time.Duration)
}

type ProviderStation struct {
	EvaNumber        int
	Code100          string
	GroupNumber      *int
	Name             string
	Lat              float32
	Lon              float32
	NoLocalTransport bool
}

type ProviderLineStop struct {
	EvaNumber int
	LineID    string
	Planned   *ProviderLineStopInfo
	Current   *ProviderLineStopInfo
	Message   string
}

type ProviderLineStopInfo struct {
	Arrival        time.Time
	Departure      time.Time
	ArrivalTrack   string
	DepartureTrack string
}

type ProviderLine struct {
	ID       string
	Type     string
	Name     string
	TripName int
	Message  string
}

type ProviderLineEdge struct {
	EvaNumberFrom        int
	EvaNumberTo          int
	LineID               string
	ProviderShortestPath *bool
	Planned              *ProviderLineStopInfo
}
