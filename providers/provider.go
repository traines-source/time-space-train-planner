package providers

import "time"

type Provider interface {
	Fetch(c Consumer)
}

type Consumer interface {
	RequestStationDataBetween(station *ProviderStation) (from time.Time, to time.Time)
	Stations() []ProviderStation
	StationByName(name string) (ProviderStation, error)
	StationByEva(evaNumber int) (ProviderStation, error)
	UpsertStation(station ProviderStation)
	UpsertLine(line ProviderLine)
	UpsertLineStop(lineStop ProviderLineStop)
}

type ProviderStation struct {
	EvaNumber   int
	Code100     string
	GroupNumber int
	Name        string
	Lat         float32
	Lon         float32
}

type ProviderLineStop struct {
	EvaNumber int
	LineID    int
	Planned   *ProviderLineStopInfo
	Current   *ProviderLineStopInfo
	Message   string
}

type ProviderLineStopInfo struct {
	Arrival   time.Time
	Departure time.Time
	Track     string
}

type ProviderLine struct {
	Type    string
	Name    string
	ID      int
	Message string
}
