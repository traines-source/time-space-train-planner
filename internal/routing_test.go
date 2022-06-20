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

	shortestPaths(map[int]*Station{0: dep, 1: interm, 2: dest}, nil, dest)

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

	shortestPaths(map[int]*Station{0: dep, 1: interm, 2: dest}, nil, dest)

	assert.Equal(t, edge1.ShortestPath, edge2, "wrong shortest path")
}

func TestRoutingTakeFirstTrainPossible(t *testing.T) {

	dep := &Station{
		Name: "Dep",
	}
	interm1 := &Station{
		Name: "Interm1",
	}
	interm2 := &Station{
		Name: "Interm2",
	}
	dest := &Station{
		Name: "Dest",
	}
	edge1 := &Edge{
		From:   dep,
		To:     interm1,
		Actual: StopInfo{Departure: minute(0), Arrival: minute(3)},
	}
	edge2 := &Edge{
		From:   interm1,
		To:     interm2,
		Actual: StopInfo{Departure: minute(4), Arrival: minute(5)},
	}
	edge3 := &Edge{
		From:   interm1,
		To:     interm2,
		Actual: StopInfo{Departure: minute(8), Arrival: minute(9)},
	}
	edge4 := &Edge{
		From:   interm2,
		To:     dest,
		Actual: StopInfo{Departure: minute(10), Arrival: minute(15)},
	}
	edge5 := &Edge{
		From:   interm1,
		To:     interm2,
		Actual: StopInfo{Departure: minute(2), Arrival: minute(5)},
	}
	edge6 := &Edge{
		From:   interm2,
		To:     interm1,
		Actual: StopInfo{Departure: minute(6), Arrival: minute(7)},
	}
	dep.Departures = []*Edge{edge1}
	interm1.Arrivals = []*Edge{edge1, edge6}
	interm1.Departures = []*Edge{edge2, edge3, edge5}
	interm2.Arrivals = []*Edge{edge2, edge3, edge5}
	interm2.Departures = []*Edge{edge4, edge6}
	dest.Arrivals = []*Edge{edge4}

	shortestPaths(map[int]*Station{0: dep, 1: interm1, 2: interm2, 3: dest}, nil, dest)

	assert.Equal(t, edge1.ShortestPath, edge2, "wrong shortest path")
	assert.Equal(t, edge2.ShortestPath, edge4, "wrong shortest path")
	assert.Equal(t, edge3.ShortestPath, edge4, "wrong shortest path")
}

func TestRoutingDontTakeFirstTrainPossibleIfTakesLonger(t *testing.T) {

	dep := &Station{
		Name: "Dep",
	}
	interm1 := &Station{
		Name: "Interm1",
	}
	interm2 := &Station{
		Name: "Interm2",
	}
	dest := &Station{
		Name: "Dest",
	}
	edge1 := &Edge{
		From:   dep,
		To:     interm1,
		Actual: StopInfo{Departure: minute(0), Arrival: minute(3)},
	}
	edge2 := &Edge{
		From:   interm1,
		To:     interm2,
		Actual: StopInfo{Departure: minute(4), Arrival: minute(9)},
	}
	edge3 := &Edge{
		From:   interm1,
		To:     interm2,
		Actual: StopInfo{Departure: minute(7), Arrival: minute(8)},
	}
	edge4 := &Edge{
		From:   interm2,
		To:     dest,
		Actual: StopInfo{Departure: minute(10), Arrival: minute(15)},
	}
	edge5 := &Edge{
		From:   interm1,
		To:     interm2,
		Actual: StopInfo{Departure: minute(2), Arrival: minute(5)},
	}
	edge6 := &Edge{
		From:   interm2,
		To:     interm1,
		Actual: StopInfo{Departure: minute(6), Arrival: minute(7)},
	}
	dep.Departures = []*Edge{edge1}
	interm1.Arrivals = []*Edge{edge1, edge6}
	interm1.Departures = []*Edge{edge2, edge3, edge5}
	interm2.Arrivals = []*Edge{edge2, edge3, edge5}
	interm2.Departures = []*Edge{edge4, edge6}
	dest.Arrivals = []*Edge{edge4}

	shortestPaths(map[int]*Station{0: dep, 1: interm1, 2: interm2, 3: dest}, nil, dest)

	assert.Equal(t, edge1.ShortestPath, edge3, "wrong shortest path 1")
	assert.Equal(t, edge2.ShortestPath, edge4, "wrong shortest path 2")
	assert.Equal(t, edge3.ShortestPath, edge4, "wrong shortest path 3")
}

func minute(minute int) time.Time {
	return time.Date(2000, 1, 1, 1, minute, 0, 0, time.Local)
}
