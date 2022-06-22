package internal

import (
	"time"
)

const inf = 1 << 31

type dijkstra struct {
	vertexAtTime dijkstraVertex
	dist         int
	hops         int
	toTarget     *Edge
}

type dijkstraToDestination struct {
	*Edge
}

type dijkstraToOrigin struct {
	*Edge
}

type dijkstraVertex interface {
	getVertices() []dijkstraVertex
	getEdge() *Edge
	isUnreachable(looseEdge *Edge) bool
	isTargetEdge(target *Edge) bool
	getShortestPath() dijkstraVertex
	setShortestPath(edge *Edge)
	travelBackDist(looseEdge *Edge) int
	earlierConnectionWithSameDist(fixedEdge *dijkstra, looseEdge *dijkstra) bool
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
	return positiveDeltaMinutes(looseEdge.Actual.Departure, edge.getEdge().Actual.Departure)
}

func (edge *dijkstraToDestination) earlierConnectionWithSameDist(fixedEdge *dijkstra, looseEdge *dijkstra) bool {
	departingEarlier := positiveDeltaMinutes(looseEdge.vertexAtTime.getEdge().Actual.Arrival, fixedEdge.vertexAtTime.getEdge().Actual.Departure) < positiveDeltaMinutes(looseEdge.vertexAtTime.getEdge().Actual.Arrival, looseEdge.toTarget.Actual.Departure)
	arrivingEarlierIfSameDestination := fixedEdge.vertexAtTime.getEdge().To != looseEdge.toTarget.To || fixedEdge.vertexAtTime.getEdge().Actual.Arrival.Before(looseEdge.toTarget.Actual.Arrival)

	return departingEarlier && arrivingEarlierIfSameDestination
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
	return positiveDeltaMinutes(edge.Actual.Arrival, looseEdge.Actual.Arrival)
}

func (edge *dijkstraToOrigin) earlierConnectionWithSameDist(fixedEdge *dijkstra, looseEdge *dijkstra) bool {
	return false
}

func shortestPaths(stations map[int]*Station, origin *Station, destination *Station) {

	for _, edgeToDestination := range destination.Arrivals {
		shortestPathsToTarget(stations, &dijkstraToDestination{edgeToDestination})
	}
	if origin != nil {
		for i := len(origin.Departures) - 1; i >= 0; i-- {
			shortestPathsToTarget(stations, &dijkstraToOrigin{origin.Departures[i]})
		}
	}
	followShortestPaths(stations)
	markEdgesAsRedundant(stations, origin, destination)
}

func shortestPathsToTarget(stations map[int]*Station, edgeToTarget dijkstraVertex) {
	verticesAtTime := buildVertexSetByTarget(edgeToTarget)

	for len(verticesAtTime) != 0 {
		u := minDist(verticesAtTime)
		u.vertexAtTime.setShortestPath(u.toTarget)
		delete(verticesAtTime, u.vertexAtTime.getEdge())

		for _, vertex := range u.vertexAtTime.getVertices() {
			if v, ok := verticesAtTime[vertex.getEdge()]; ok {
				alt := u.dist + u.vertexAtTime.travelBackDist(v.vertexAtTime.getEdge())
				if alt < v.dist || alt == v.dist && u.hops+1 <= v.hops && u.vertexAtTime.earlierConnectionWithSameDist(u, v) {
					v.dist = alt
					v.hops = u.hops + 1
					v.toTarget = u.vertexAtTime.getEdge()
				}
			}
		}
	}
}

func buildVertexSetByTarget(edgeToTarget dijkstraVertex) map[*Edge]*dijkstra {
	verticesAtTime := map[*Edge]*dijkstra{}
	verticesAtTime[edgeToTarget.getEdge()] = &dijkstra{
		vertexAtTime: edgeToTarget,
		dist:         int(edgeToTarget.getEdge().Actual.Arrival.Sub(edgeToTarget.getEdge().Actual.Departure).Minutes()),
		hops:         0,
		toTarget:     nil,
	}
	buildVertexSet(verticesAtTime, edgeToTarget, edgeToTarget.getEdge())
	return verticesAtTime
}

func buildVertexSet(verticesAtTime map[*Edge]*dijkstra, vertexAtTime dijkstraVertex, targetEdge *Edge) {
	for _, edge := range vertexAtTime.getVertices() {
		if edge.getShortestPath() != nil {
			continue
		}
		if vertexAtTime.isUnreachable(edge.getEdge()) {
			continue
		}
		if edge.isTargetEdge(targetEdge) {
			continue
		}
		if _, ok := verticesAtTime[edge.getEdge()]; ok {
			continue
		}
		verticesAtTime[edge.getEdge()] = &dijkstra{
			vertexAtTime: edge,
			dist:         inf,
			hops:         inf,
			toTarget:     nil,
		}
		buildVertexSet(verticesAtTime, edge, targetEdge)
	}
}

func minDist(verticesAtTime map[*Edge]*dijkstra) *dijkstra {
	var minVertex *dijkstra
	for _, vertex := range verticesAtTime {
		if minVertex == nil || vertex.dist < minVertex.dist {
			minVertex = vertex
		}
	}
	return minVertex
}

func positiveDeltaMinutes(from time.Time, to time.Time) int {
	min := int(to.Sub(from).Minutes())
	if min < 0 {
		return inf
	}
	return min
}

func markEdgesAsRedundant(stations map[int]*Station, origin *Station, destination *Station) {
	for _, station := range stations {
		for _, departure := range station.Departures {
			markEdgeAsRedundant(departure, origin, destination)
			markEdgeAsDiscarded(departure)
		}
	}
}

func markEdgeAsRedundant(edge *Edge, origin *Station, destination *Station) {
	if edge.ShortestPath == nil && edge.To != destination || edge.ReverseShortestPath == nil && edge.From != origin {
		edge.Redundant = true
		return
	}
	if revisitsSameStation(edge) {
		edge.Redundant = true
		return
	}
	if len(edge.ShortestPathFor) > 1 {
		edge.Redundant = false
		return
	} else if edge.Line.Type == "Foot" {
		edge.Redundant = true
		return
	}
	if existsLaterDepartureWithEarlierArrival(edge, destination) {
		edge.Redundant = true
		return
	}
}

func markEdgeAsDiscarded(edge *Edge) {
	if edge.Redundant && edge.From.GroupNumber != nil && edge.To.GroupNumber != nil && *edge.From.GroupNumber == *edge.To.GroupNumber {
		edge.Discarded = true
	}
}

func existsLaterDepartureWithEarlierArrival(edge *Edge, destination *Station) bool {
	hopsByEdge, arrivalByEdge := calculateTravelLength(edge, destination)
	if hopsByEdge == -1 {
		return true
	}

	for _, departure := range edge.From.Departures {
		if departure == edge || departure.Actual.Departure.Before(edge.Actual.Departure) {
			continue
		}
		hops, arrival := calculateTravelLength(departure, destination)
		if hops != -1 && (arrival.Before(arrivalByEdge)) {
			return true
		}
	}
	return false
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

func revisitsSameStation(edge *Edge) bool {
	from := edge.From.EvaNumber
	nextEdge := edge
	for {
		if nextEdge.ShortestPath != nil {
			nextEdge = nextEdge.ShortestPath
			if nextEdge.To.EvaNumber == from {
				return true
			}
		} else {
			return false
		}
	}
}

func followShortestPaths(stations map[int]*Station) {
	for _, l := range stations {
		for _, origin := range l.Departures {
			for e := origin; e != nil; e = e.ShortestPath {
				setShortestPathFor(origin, e)
			}
			for e := origin; e != nil; e = e.ReverseShortestPath {
				setShortestPathFor(origin, e)
			}
		}
	}
}

func setShortestPathFor(origin *Edge, e *Edge) {
	if e.ShortestPathFor == nil {
		e.ShortestPathFor = map[*Edge]struct{}{}
	}
	e.ShortestPathFor[origin] = struct{}{}
}
