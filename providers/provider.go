package providers

import "time"

type Provider interface {
	Vias(c Consumer) error
	DeparturesArrivals(c Consumer) error
	Enrich(c Consumer) error
}

type Consumer interface {
	RequestStationDataBetween(station *ProviderStation) (from time.Time, to time.Time)
	Stations() []ProviderStation
	StationByName(name string) (ProviderStation, error)
	StationByID(id string) (ProviderStation, error)
	RegionalOnly() bool
	UpsertStation(station ProviderStation)
	UpsertLine(line ProviderLine)
	UpsertLineStop(lineStop ProviderLineStop)
	UpsertLineEdge(lineEdge ProviderLineEdge)
	SetExpectedTravelDuration(duration time.Duration)
}

type ProviderStation struct {
	ID               string
	Code100          string
	GroupID          *string
	AltID            *string
	Name             string
	Lat              float64
	Lon              float64
	NoLocalTransport bool
}

type ProviderLineStop struct {
	ID        string
	LineID    string
	Planned   *ProviderLineStopInfo
	Current   *ProviderLineStopInfo
	Message   string
	Cancelled bool
}

type ProviderLineStopInfo struct {
	Arrival        time.Time
	Departure      time.Time
	ArrivalTrack   string
	DepartureTrack string
}

type ProviderLine struct {
	ID        string
	Type      string
	Name      string
	TripName  int
	Message   string
	Direction string
}

type ProviderLineEdge struct {
	IDFrom               string
	IDTo                 string
	LineID               string
	ProviderShortestPath *bool
	Planned              *ProviderLineStopInfo
}
