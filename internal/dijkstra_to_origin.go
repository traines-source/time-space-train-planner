package internal

type dijkstraToOrigin struct {
	*Edge
}

func (edge *dijkstraToOrigin) getVertices() []dijkstraVertex {
	var edges []dijkstraVertex
	for _, edge := range edge.To.Departures {
		edges = append(edges, &dijkstraToOrigin{edge})
	}
	return edges
}

func (edge *dijkstraToOrigin) getEdge() *Edge {
	return edge.Edge
}

func (edge *dijkstraToOrigin) isUnreachable(looseEdge *Edge) bool {
	return looseEdge.Actual.Departure.Before(edge.getEdge().Actual.Arrival)
}

func (edge *dijkstraToOrigin) isTargetEdge(targetEdge *Edge) bool {
	return edge.From == targetEdge.From
}

func (edge *dijkstraToOrigin) getShortestPath() dijkstraVertex {
	if edge.Edge.ReverseShortestPath == nil {
		return nil
	}
	return &dijkstraToOrigin{edge.Edge.ReverseShortestPath}
}

func (edge *dijkstraToOrigin) setShortestPath(target *Edge) {
	edge.Edge.ReverseShortestPath = target
}

func (edge *dijkstraToOrigin) travelBackDist(looseEdge *Edge) int {
	return deltaMinutes(edge.Actual.Arrival, looseEdge.Actual.Arrival)
}

func (edge *dijkstraToOrigin) earlierConnectionWithSameDist(fixedEdge *dijkstra, looseEdge *dijkstra) bool {
	return false
}