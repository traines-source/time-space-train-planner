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

type dijkstraVertex interface {
	getVertices() []dijkstraVertex
	getEdge() *Edge
	isUnreachable(looseEdge *Edge) bool
	isTargetEdge(target *Edge) bool
	getShortestPath() dijkstraVertex
	setShortestPath(edge *Edge)
	setTargetTime(minutes int)
	travelBackDist(looseEdge *Edge) int
	earlierConnectionWithSameDist(fixedEdge *dijkstra, looseEdge *dijkstra) bool
}

func shortestPaths(stations map[int]*Station, origin *Station, destination *Station, regionly bool) {

	for _, edgeToDestination := range destination.Arrivals {
		shortestPathsToTarget(stations, &dijkstraToDestination{edgeToDestination}, regionly)
	}
	if origin != nil {
		for i := len(origin.Departures) - 1; i >= 0; i-- {
			shortestPathsToTarget(stations, &dijkstraToOrigin{origin.Departures[i]}, regionly)
		}
	}
	followShortestPaths(stations)
	markEdgesAsRedundant(stations, origin, destination, regionly)
}

func shortestPathsToTarget(stations map[int]*Station, edgeToTarget dijkstraVertex, regionly bool) {
	if isNotEligible(edgeToTarget.getEdge(), regionly) {
		return
	}
	verticesAtTime := buildVertexSetByTarget(edgeToTarget, regionly)

	for len(verticesAtTime) != 0 {
		u := minDist(verticesAtTime)
		u.vertexAtTime.setShortestPath(u.toTarget)
		u.vertexAtTime.setTargetTime(u.dist)
		delete(verticesAtTime, u.vertexAtTime.getEdge())

		for _, vertex := range u.vertexAtTime.getVertices() {
			if v, ok := verticesAtTime[vertex.getEdge()]; ok {
				travelBack := u.vertexAtTime.travelBackDist(v.vertexAtTime.getEdge())
				if travelBack < 0 {
					continue
				}
				alt := u.dist + travelBack
				if alt < v.dist || alt == v.dist && u.hops+1 <= v.hops && u.vertexAtTime.earlierConnectionWithSameDist(u, v) {
					v.dist = alt
					v.hops = u.hops + 1
					v.toTarget = u.vertexAtTime.getEdge()
				}
			}
		}
	}
}

func isNotEligible(e *Edge, regionly bool) bool {
	return e.Cancelled || regionly && (e.Line.Type == "national" || e.Line.Type == "nationalExpress")
}

func buildVertexSetByTarget(edgeToTarget dijkstraVertex, regionly bool) map[*Edge]*dijkstra {
	verticesAtTime := map[*Edge]*dijkstra{}
	verticesAtTime[edgeToTarget.getEdge()] = &dijkstra{
		vertexAtTime: edgeToTarget,
		dist:         int(edgeToTarget.getEdge().Actual.Arrival.Sub(edgeToTarget.getEdge().Actual.Departure).Minutes()),
		hops:         0,
		toTarget:     nil,
	}
	buildVertexSet(verticesAtTime, edgeToTarget, edgeToTarget.getEdge(), regionly)
	return verticesAtTime
}

func buildVertexSet(verticesAtTime map[*Edge]*dijkstra, vertexAtTime dijkstraVertex, targetEdge *Edge, regionly bool) {
	for _, edge := range vertexAtTime.getVertices() {
		if edge.getShortestPath() != nil {
			continue
		}
		if vertexAtTime.isUnreachable(edge.getEdge()) {
			continue
		}
		if isNotEligible(edge.getEdge(), regionly) {
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
		buildVertexSet(verticesAtTime, edge, targetEdge, regionly)
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

func deltaMinutes(from time.Time, to time.Time) int {
	return int(to.Sub(from).Minutes())
}

func markEdgesAsRedundant(stations map[int]*Station, origin *Station, destination *Station, regionly bool) {
	for _, station := range stations {
		for _, departure := range station.Departures {
			markEdgeAsRedundant(departure, origin, destination, regionly)
			markEdgeAsDiscarded(departure)
		}
	}
}

func markEdgeAsRedundant(edge *Edge, origin *Station, destination *Station, regionly bool) {
	if edge.ShortestPath == nil && edge.To != destination || edge.ReverseShortestPath == nil && edge.From != origin {
		edge.Redundant = true
		return
	}
	if isNotEligible(edge, regionly) {
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
	if edge.Redundant && edge.Line.Type == "Foot" {
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
