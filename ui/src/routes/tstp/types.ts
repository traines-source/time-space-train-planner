export type Edge = any;
export type Coord = any;
export type Station = any;
export class Response {
    DefaultShortestPathID: string | undefined;
    Edges: Edge[string] = {};
    From: Station;
    To: Station;
    SortedEdges: string[] = [];
    Stations: Station[string] = [];
    TimeIndicators: Coord = [];
    constructor(public MaxSpace: number, public MaxTime: string, public MinTime: string, public SpaceAxisSize: number, public TimeAxisSize: number, public TimeAxisDistance: number) {
    }
}

export class Selection {
    edge: Edge | undefined;
    station: Station | undefined;
    from: Date | undefined;
    to: Date | undefined;

    static fromEdge(edge: Edge) {
        const t = new Selection();
        t.edge = edge;
        t.station = undefined;
        t.from = undefined;
        t.to = undefined;
        return t;
    }

    static fromStation(station: Station) {
        const t = new Selection();
        t.edge = undefined;
        t.station = station;
        t.from = undefined;
        t.to = undefined;
        return t;
    }
}
