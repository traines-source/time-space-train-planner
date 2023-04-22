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

