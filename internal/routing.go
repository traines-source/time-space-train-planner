package internal

import (
	"log"
)

const inf = 1 << 31

type dijkstra struct {
	vertexAtDeparture *Edge
	dist              int
	previous          *Edge
}

func shortestPaths(stations map[int]*Station, destination *Station) {

	for _, edgeToDestination := range destination.Arrivals {
		shortestPathsToEdge(stations, edgeToDestination)
	}
	markAsRedundantIfNoShortestPath(stations, destination)
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
				if alt < v.dist {
					// || alt == v.dist && v.previous.Actual.Departure.Sub(v.vertexAtDeparture.Actual.Arrival) > u.vertexAtDeparture.Actual.Departure.Sub(v.vertexAtDeparture.Actual.Arrival)
					v.dist = alt
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
		verticesAtDeparture[edge] = &dijkstra{
			vertexAtDeparture: edge,
			dist:              inf,
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
	min := int(next.Actual.Departure.Sub(previous.Actual.Departure).Minutes())
	if min <= 0 {
		return inf
	}
	return min
}

func markAsRedundantIfRevisitsSameStation(edge *dijkstra) {
	from := edge.vertexAtDeparture.From.EvaNumber
	nextEdge := edge.vertexAtDeparture
	for {
		if nextEdge.ShortestPath != nil {
			nextEdge = nextEdge.ShortestPath
			if nextEdge.To.EvaNumber == from {
				edge.vertexAtDeparture.Redundant = true
				log.Print("redundant")
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
