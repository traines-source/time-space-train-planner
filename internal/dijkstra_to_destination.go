package internal

type dijkstraToDestination struct {
	*Edge
}

func (edge *dijkstraToDestination) getVertices() []dijkstraVertex {
	var edges []dijkstraVertex
	for _, edge := range edge.From.Arrivals {
		edges = append(edges, &dijkstraToDestination{edge})
	}
	for i, j := 0, len(edges)-1; i < j; i, j = i+1, j-1 {
		edges[i], edges[j] = edges[j], edges[i]
	}
	return edges
}

func (edge *dijkstraToDestination) getEdge() *Edge {
	return edge.Edge
}

func (edge *dijkstraToDestination) isUnreachable(looseEdge *Edge) bool {
	return edge.getEdge().Actual.Departure.Before(looseEdge.Actual.Arrival)
}

func (edge *dijkstraToDestination) isTargetEdge(targetEdge *Edge) bool {
	return edge.To == targetEdge.To
}

func (edge *dijkstraToDestination) getShortestPath() dijkstraVertex {
	if edge.Edge.ShortestPath == nil {
		return nil
	}
	return &dijkstraToDestination{edge.Edge.ShortestPath}
}

func (edge *dijkstraToDestination) setShortestPath(target *Edge) {
	edge.Edge.ShortestPath = target
}

func (edge *dijkstraToDestination) travelBackDist(looseEdge *Edge) int {
	return deltaMinutes(looseEdge.Actual.Departure, edge.getEdge().Actual.Departure)
}

func (edge *dijkstraToDestination) earlierConnectionWithSameDist(fixedEdge *dijkstra, looseEdge *dijkstra) bool {
	a := deltaMinutes(looseEdge.vertexAtTime.getEdge().Actual.Arrival, fixedEdge.vertexAtTime.getEdge().Actual.Departure)
	b := deltaMinutes(looseEdge.vertexAtTime.getEdge().Actual.Arrival, looseEdge.toTarget.Actual.Departure)
	if a < 0 || b < 0 {
		return false
	}
	departingEarlier := a < b
	arrivingEarlierIfSameDestination := fixedEdge.vertexAtTime.getEdge().To != looseEdge.toTarget.To || fixedEdge.vertexAtTime.getEdge().Actual.Arrival.Before(looseEdge.toTarget.Actual.Arrival)

	return departingEarlier && arrivingEarlierIfSameDestination
}