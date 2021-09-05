package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleRouting(t *testing.T) {

	dep := &Station{
		Name: "Dep",
	}
	interm := &Station{
		Name: "Interm",
	}
	dest := &Station{
		Name: "Dest",
	}
	edge1 := &Edge{
		From:   dep,
		To:     interm,
		Actual: StopInfo{Departure: minute(1), Arrival: minute(3)},
	}
	edge2 := &Edge{
		From:   interm,
		To:     dest,
		Actual: StopInfo{Departure: minute(5), Arrival: minute(10)},
	}
	edge3 := &Edge{
		From:   dep,
		To:     interm,
		Actual: StopInfo{Departure: minute(3), Arrival: minute(4)},
	}
	edge4 := &Edge{
		From:   interm,
		To:     dest,
		Actual: StopInfo{Departure: minute(20), Arrival: minute(22)},
	}
	edge5 := &Edge{
		From:   dep,
		To:     interm,
		Actual: StopInfo{Departure: minute(3), Arrival: minute(6)},
	}
	dep.Departures = []*Edge{edge1, edge3, edge5}
	interm.Arrivals = []*Edge{edge1, edge3, edge5}
	interm.Departures = []*Edge{edge2, edge4}
	dest.Arrivals = []*Edge{edge2, edge4}

	shortestPaths(map[int]*Station{0: dep, 1: interm, 2: dest}, dest)

	assert.Equal(t, edge1.ShortestPath, edge2, "wrong shortest path")
	assert.Equal(t, edge3.ShortestPath, edge2, "wrong shortest path")
	assert.Equal(t, edge5.ShortestPath, edge4, "wrong shortest path")
}

func TestRoutingWithOvertakingTrains(t *testing.T) {

	dep := &Station{
		Name: "Dep",
	}
	interm := &Station{
		Name: "Interm",
	}
	dest := &Station{
		Name: "Dest",
	}
	edge1 := &Edge{
		From:   dep,
		To:     interm,
		Actual: StopInfo{Departure: minute(0), Arrival: minute(3)},
	}
	edge2 := &Edge{
		From:   interm,
		To:     dest,
		Actual: StopInfo{Departure: minute(6), Arrival: minute(10)},
	}
	edge4 := &Edge{
		From:   interm,
		To:     dest,
		Actual: StopInfo{Departure: minute(4), Arrival: minute(15)},
	}
	dep.Departures = []*Edge{edge1}
	interm.Arrivals = []*Edge{edge1}
	interm.Departures = []*Edge{edge2, edge4}
	dest.Arrivals = []*Edge{edge2, edge4}

	shortestPaths(map[int]*Station{0: dep, 1: interm, 2: dest}, dest)

	assert.Equal(t, edge1.ShortestPath, edge2, "wrong shortest path")
}

func minute(minute int) time.Time {
	return time.Date(2000, 1, 1, 1, minute, 0, 0, time.Local)
}
