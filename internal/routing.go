package internal

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
}

func shortestPathsToEdge(stations map[int]*Station, edgeToDestination *Edge) {
	verticesAtDeparture := buildVertexSetByDestination(edgeToDestination)

	for len(verticesAtDeparture) != 0 {
		u := minDist(verticesAtDeparture)
		u.vertexAtDeparture.ShortestPath = u.previous
		delete(verticesAtDeparture, u.vertexAtDeparture)

		for _, vertex := range u.vertexAtDeparture.From.Arrivals {
			if v, ok := verticesAtDeparture[vertex]; ok {
				alt := u.dist + travelBackDist(u.vertexAtDeparture, v.vertexAtDeparture)
				if alt < v.dist {
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
	buildVertexSet(verticesAtDeparture, edgeToDestination)
	return verticesAtDeparture
}

func buildVertexSet(verticesAtDeparture map[*Edge]*dijkstra, vertexAtDeparture *Edge) {
	for _, edge := range vertexAtDeparture.From.Arrivals {
		if edge.ShortestPath != nil {
			continue
		}
		if vertexAtDeparture.Actual.Departure.Before(edge.Actual.Arrival) {
			continue
		}
		verticesAtDeparture[edge] = &dijkstra{
			vertexAtDeparture: edge,
			dist:              inf,
			previous:          nil,
		}
		buildVertexSet(verticesAtDeparture, edge)
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
	return int(next.Actual.Departure.Sub(previous.Actual.Departure).Minutes())
}
