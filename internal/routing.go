package internal

import (
	"time"
)

const inf = 1 << 31

type dijkstra struct {
	vertexAtDeparture *Edge
	dist              int
	hops              int
	previous          *Edge
}

func shortestPaths(stations map[int]*Station, destination *Station) {

	for _, edgeToDestination := range destination.Arrivals {
		shortestPathsToEdge(stations, edgeToDestination)
	}
	//markAsRedundantIfNoShortestPath(stations, destination)
	markEdgesAsRedundant(stations, destination)
}

func shortestPathsToEdge(stations map[int]*Station, edgeToDestination *Edge) {
	verticesAtDeparture := buildVertexSetByDestination(edgeToDestination)

	for len(verticesAtDeparture) != 0 {
		u := minDist(verticesAtDeparture)
		u.vertexAtDeparture.ShortestPath = u.previous
		markAsRedundantIfRevisitsSameStation(u)
		delete(verticesAtDeparture, u.vertexAtDeparture)

		for _, vertex := range u.vertexAtDeparture.From.Arrivals {
			if v, ok := verticesAtDeparture[vertex]; ok {
				alt := u.dist + travelBackDist(u.vertexAtDeparture, v.vertexAtDeparture)
				if alt < v.dist || alt == v.dist && u.hops+1 <= v.hops && earlierConnectionWithSameDist(u, v) {
					v.dist = alt
					v.hops = u.hops + 1
					v.previous = u.vertexAtDeparture
				}
			}
		}
	}
}

func buildVertexSetByDestination(edgeToDestination *Edge) map[*Edge]*dijkstra {
	verticesAtDeparture := map[*Edge]*dijkstra{}
	verticesAtDeparture[edgeToDestination] = &dijkstra{
		vertexAtDeparture: edgeToDestination,
		dist:              int(edgeToDestination.Actual.Arrival.Sub(edgeToDestination.Actual.Departure).Minutes()),
		hops:              0,
		previous:          nil,
	}
	buildVertexSet(verticesAtDeparture, edgeToDestination, edgeToDestination.To)
	return verticesAtDeparture
}

func buildVertexSet(verticesAtDeparture map[*Edge]*dijkstra, vertexAtDeparture *Edge, destination *Station) {
	for _, edge := range vertexAtDeparture.From.Arrivals {
		if edge.ShortestPath != nil {
			continue
		}
		if vertexAtDeparture.Actual.Departure.Before(edge.Actual.Arrival) {
			continue
		}
		if edge.To == destination {
			continue
		}
		if _, ok := verticesAtDeparture[edge]; ok {
			continue
		}
		verticesAtDeparture[edge] = &dijkstra{
			vertexAtDeparture: edge,
			dist:              inf,
			hops:              inf,
			previous:          nil,
		}
		buildVertexSet(verticesAtDeparture, edge, destination)
	}
}

func minDist(verticesAtDeparture map[*Edge]*dijkstra) *dijkstra {
	var minVertex *dijkstra
	for _, vertex := range verticesAtDeparture {
		if minVertex == nil || vertex.dist < minVertex.dist {
			minVertex = vertex
		}
	}
	return minVertex
}

func travelBackDist(next *Edge, previous *Edge) int {
	return positiveDeltaMinutes(previous.Actual.Departure, next.Actual.Departure)
}

func positiveDeltaMinutes(previous time.Time, next time.Time) int {
	min := int(next.Sub(previous).Minutes())
	if min < 0 {
		return inf
	}
	return min
}

func earlierConnectionWithSameDist(u *dijkstra, v *dijkstra) bool {
	return positiveDeltaMinutes(v.vertexAtDeparture.Actual.Arrival, u.vertexAtDeparture.Actual.Departure) < positiveDeltaMinutes(v.vertexAtDeparture.Actual.Arrival, v.previous.Actual.Departure)
}

func markEdgesAsRedundant(stations map[int]*Station, destination *Station) {
	for _, station := range stations {
		for _, departure := range station.Departures {
			markEdgeAsRedundant(departure, destination)
		}
	}
}

func markEdgeAsRedundant(edge *Edge, destination *Station) {
	hopsByEdge, arrivalByEdge := calculateTravelLength(edge, destination)
	if hopsByEdge == -1 {
		edge.Redundant = true
		return
	}

	for _, departure := range edge.From.Departures {
		if departure == edge || departure.Actual.Departure.Before(edge.Actual.Departure) {
			continue
		}
		hops, arrival := calculateTravelLength(departure, destination)
		if hops != -1 && (arrival.Before(arrivalByEdge)) {
			edge.Redundant = true
			break
		}
	}
}

func calculateTravelLength(startEdge *Edge, destination *Station) (hops int32, arrival time.Time) {
	nextEdge := startEdge
	for {
		if nextEdge.To == destination {
			return hops, nextEdge.Actual.Arrival
		}
		if nextEdge.ShortestPath != nil {
			nextEdge = nextEdge.ShortestPath
		} else {
			break
		}
		hops++
	}
	return -1, time.Time{}
}

func markAsRedundantIfRevisitsSameStation(edge *dijkstra) {
	from := edge.vertexAtDeparture.From.EvaNumber
	nextEdge := edge.vertexAtDeparture
	for {
		if nextEdge.ShortestPath != nil {
			nextEdge = nextEdge.ShortestPath
			if nextEdge.To.EvaNumber == from {
				edge.vertexAtDeparture.Redundant = true
				break
			}
		} else {
			break
		}
	}
}

func markAsRedundantIfNoShortestPath(stations map[int]*Station, destination *Station) {
	for _, station := range stations {
		for _, departure := range station.Departures {
			if departure.ShortestPath == nil && departure.To != destination {
				departure.Redundant = true
			}
		}
	}
}
