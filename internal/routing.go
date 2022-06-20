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
	isTargetEdge(*Edge) bool
	getShortestPath() dijkstraVertex
	setShortestPath(*Edge)
	earlierConnectionWithSameDist(*dijkstra, *dijkstra) bool
}

func (edge *dijkstraToDestination) getVertices() []dijkstraVertex {
	var edges []dijkstraVertex
	for _, edge := range edge.From.Arrivals {
		edges = append(edges, &dijkstraToDestination{edge})
	}
	return edges
}

func (edge *dijkstraToDestination) getEdge() *Edge {
	return edge.Edge
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

func (edge *dijkstraToDestination) earlierConnectionWithSameDist(u *dijkstra, v *dijkstra) bool {
	departingEarlier := positiveDeltaMinutes(v.vertexAtTime.getEdge().Actual.Arrival, u.vertexAtTime.getEdge().Actual.Departure) < positiveDeltaMinutes(v.vertexAtTime.getEdge().Actual.Arrival, v.toTarget.Actual.Departure)
	arrivingEarlierIfSameDestination := u.vertexAtTime.getEdge().To != v.toTarget.To || u.vertexAtTime.getEdge().Actual.Arrival.Before(v.toTarget.Actual.Arrival)

	return departingEarlier && arrivingEarlierIfSameDestination
}

/*
func TODOearlierConnectionWithSameDist(u *dijkstra, v *dijkstra) bool {
	arrivingLater := positiveDeltaMinutes(u.vertexAtTime.Actual.Arrival, v.vertexAtTime.Actual.Departure) < positiveDeltaMinutes(v.toTarget.Actual.Arrival, v.vertexAtTime.Actual.Departure)
	departingLaterIfSameOrigin := u.vertexAtTime.From != v.toTarget.From || v.toTarget.Actual.Departure.Before(u.vertexAtTime.Actual.Departure)

	return arrivingLater && departingLaterIfSameOrigin
}*/

func shortestPaths(stations map[int]*Station, destination *Station) {

	for _, edgeToDestination := range destination.Arrivals {
		dE := &dijkstraToDestination{edgeToDestination}
		shortestPathsToDestination(stations, dE)
	}
	//markAsRedundantIfNoShortestPath(stations, destination)
	markEdgesAsRedundant(stations, destination)
}

func shortestPathsToDestination(stations map[int]*Station, edgeToTarget dijkstraVertex) {
	verticesAtDeparture := buildVertexSetByTarget(edgeToTarget)

	for len(verticesAtDeparture) != 0 {
		u := minDist(verticesAtDeparture)
		u.vertexAtTime.setShortestPath(u.toTarget)
		markAsRedundantIfRevisitsSameStation(u)
		delete(verticesAtDeparture, u.vertexAtTime.getEdge())

		for _, vertex := range u.vertexAtTime.getVertices() {
			if v, ok := verticesAtDeparture[vertex.getEdge()]; ok {
				alt := u.dist + travelBackDist(u.vertexAtTime.getEdge(), v.vertexAtTime.getEdge())
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
		if vertexAtTime.getEdge().Actual.Departure.Before(edge.getEdge().Actual.Arrival) {
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

func travelBackDist(fixedEdge *Edge, looseEdge *Edge) int {
	return positiveDeltaMinutes(looseEdge.Actual.Departure, fixedEdge.Actual.Departure)
}

func travelForwardDist(fixedEdge *Edge, looseEdge *Edge) int {
	return positiveDeltaMinutes(fixedEdge.Actual.Arrival, looseEdge.Actual.Arrival)
}

func positiveDeltaMinutes(from time.Time, to time.Time) int {
	min := int(to.Sub(from).Minutes())
	if min < 0 {
		return inf
	}
	return min
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
	from := edge.vertexAtTime.getEdge().From.EvaNumber
	nextEdge := edge.vertexAtTime
	for {
		if nextEdge.getShortestPath() != nil {
			nextEdge = nextEdge.getShortestPath()
			if nextEdge.getEdge().To.EvaNumber == from {
				edge.vertexAtTime.getEdge().Redundant = true
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
