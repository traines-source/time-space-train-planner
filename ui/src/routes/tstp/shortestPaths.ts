import { parseTime } from "./labels";
import type { Edge, Station, Response } from "./types";

function getStationsInGroup(data: Response, station: Station): Station[] {
    const relevantStations = [];
    for (let s of Object.values(data.Stations)) {
        if (s.ID == station.ID || s.GroupID == station.GroupID) {
            relevantStations.push(s);
        }
    }
    return relevantStations;
}

function hasDistribution(e: Edge) {
    return e.DestinationArrival && e.DestinationArrival.Histogram && e.DestinationArrival.Histogram.length;
}

function calcNextDepartureIndex(station: Station, relevantStations: Station[], indices: number[], earliestDepartureTime: (e: Edge) => number, edgeResolver: (id: string) => Edge, stationResolver: (id: string) => Station): number | undefined {
    let nextDepartureIndex = undefined;
    for (let s=0; s<relevantStations.length; s++) {
        let e;
        while(true) {
            if (indices[s] >= relevantStations[s].BestDepartures.length) {
                e = undefined;
                break;
            }
            e = edgeResolver(relevantStations[s].BestDepartures[indices[s]]);
            let departure = parseTime(e.Actual.Departure);
            if (departure >= earliestDepartureTime(e)
            && !(hasDistribution(e) && e.DestinationArrival.FeasibleProbability < 0.1)
            && !(e.Line?.Type == "Foot" && station.ID != relevantStations[s].ID)
            && !(stationResolver(e.To.SpaceAxis).GroupID == station.GroupID)) {
                break;
            }
            indices[s]++;
        }
        if (!e) {
            continue;
        }
        if (nextDepartureIndex == undefined) {
            nextDepartureIndex = s;    
        } else {
            const nextDeparture = edgeResolver(relevantStations[nextDepartureIndex].BestDepartures[indices[nextDepartureIndex]]);
            if ((nextDeparture.DestinationArrival?.Mean || nextDeparture.EarliestDestinationArrival) > (e.DestinationArrival?.Mean || e.EarliestDestinationArrival)) {
                nextDepartureIndex = s;
            }
        }
    }
    return nextDepartureIndex;
}

function walkingDistance(fromId: string, toId: string, stationResolver: (id: string) => Station): number {
    const from = stationResolver(fromId);
    const to = stationResolver(toId);
    const φ1 = from.Lat * Math.PI/180, φ2 = to.Lat * Math.PI/180, Δλ = (to.Lon-from.Lon) * Math.PI/180, R = 6371e3;
    return Math.acos( Math.sin(φ1)*Math.sin(φ2) + Math.cos(φ1)*Math.cos(φ2) * Math.cos(Δλ) )*R;
}

function walkingDurationMs(fromId: string, toId: string, stationResolver: (id: string) => Station): number {
    return walkingDistance(fromId, toId, stationResolver)/5*3600;
}

export {
    hasDistribution, calcNextDepartureIndex, walkingDistance, walkingDurationMs, getStationsInGroup
}