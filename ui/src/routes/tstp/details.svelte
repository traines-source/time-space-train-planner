<script lang="ts">
    import { t } from '$lib/translations';
    import { store } from "../store"
    import { optionsQueryString } from "../query"
    import Footer from '../footer.svelte'; 
    import {label, type, departure, arrival, parseTime, simpleTime} from './labels';
    import type { Edge, Station, Response } from './types';

    export let selection: Edge;
    export let loading: boolean;
    export let doRefresh: () => void;
    export let selectEdge: (id: string) => void;
    export let selectStation: (id: string) => void;
    export let data: Response;
    export let error: string | undefined;

    let selectedEdgeHistory: Edge[] = [];

    let nextBestDepartures: Edge[] | undefined = undefined;
    let nextBestDeparturesBoundsMinutes: number[] = [0,0];
    const maxNextBestDepartures = 5;
    const negativeTransferMinutes = 5;
    $: {
        if (selection.edge) {
            nextBestDeparturesForEdge();
            rectifyEdgeHistory();
        } else if (selection.station) {
            nextBestDeparturesForStation();
        }
    }

    function nextBestDeparturesForEdge() {
        if (selection.edge.To.SpaceAxis == data.To.ID) {
            nextBestDepartures = undefined;
            return;
        }
        const station = stationResolver(selection.edge.To.SpaceAxis);
        if (!selection.from) {
            selection.from = new Date(parseTime(selection.edge.Actual.Arrival)-negativeTransferMinutes*60*1000);
        }
        updateNextBestDepartures(station, selection.from);
    }

    function nextBestDeparturesForStation() {
        if (!selection.from) {
            let n = new Date().getTime();
            if (selectedEdgeHistory.length > 0) {
                n = parseTime(selectedEdgeHistory[selectedEdgeHistory.length-1].Actual.Departure);
            } else if (n < new Date(data.MinTime).getTime() || n > new Date(data.MaxTime).getTime()) {
                n = new Date(data.MinTime).getTime();
            }            
            selection.from = new Date(n-negativeTransferMinutes*60*1000);
        }
        selectedEdgeHistory = [];
        updateNextBestDepartures(selection.station, selection.from);
    }

    function getBounds(e: Edge): number[] {
        const h = e.DestinationArrival.Histogram;
        const start = parseTime(e.DestinationArrival.Start)/60/1000;
        let lower = undefined;
        let upper = undefined;
        let lowerAccum = 0;
        let upperAccum = 0;
        let thresh = 0.05;
        for (let i=0;i<h.length;i++) {
            lowerAccum += h[i];
            upperAccum += h[h.length-1-i];
            if (lower == undefined && lowerAccum > thresh) lower = start+i-1;
            if (upper == undefined && upperAccum > thresh) upper = start+h.length-i;
            if (lower && upper) break;
        }
        return [lower, upper];
    }

    function updateNextBestDepartures(station: Station, time: Date) {        
        const candidates = [];
        
        let lowerBound = undefined;
        let upperBound = undefined;
        for (let i=0; i<station.BestDepartures.length; i++) {
            if (candidates.length >= maxNextBestDepartures) {
                break;
            }
            const e = edgeResolver(station.BestDepartures[i]);
            if (parseTime(e.Actual.Departure) < time.getTime()) {
                continue;
            }
            candidates.push(e);
            if (!hasDistribution(e)) continue;
            const bounds = getBounds(e);
            if (lowerBound == undefined || bounds[0] < lowerBound) {
                lowerBound = bounds[0];
            }
            if (upperBound == undefined || bounds[1] > upperBound) {
                upperBound = bounds[1];
            }            
        }
        candidates.sort((a, b) => {
            return parseTime(a.Planned.Departure)-parseTime(b.Planned.Departure);
        });
        
        const drawWidth = 100;
        const padding = (upperBound-lowerBound)/(drawWidth/20-1);
        nextBestDeparturesBoundsMinutes = [lowerBound-padding, upperBound+padding];
        nextBestDepartures = candidates;
    }

    function rectifyEdgeHistory() {
        if (selectedEdgeHistory[selectedEdgeHistory.length-1] != selection.edge) {
            selectedEdgeHistory = [selection.edge];
            console.log("history reset");
        }
    }

    function edgeResolver(id: string): Edge {
        if (!id) return undefined;
        return data.Edges[id];
    }

    function stationResolver(id: string): Station {
        if (!id) return undefined;
        return data.Stations[id];
    }

    function hasDistribution(e: Edge) {
        return e.DestinationArrival && e.DestinationArrival.Histogram && e.DestinationArrival.Histogram.length;
    }

    function getDrawXRatio() {
        const drawWidth = 100;
        return (nextBestDeparturesBoundsMinutes[1]-nextBestDeparturesBoundsMinutes[0])/drawWidth;
    }

    function histogram(e: Edge) {
        console.log(e.Line.Name, e.DestinationArrival.Start, e.DestinationArrival.Mean, e.DestinationArrival.Histogram.map(p => Math.round(p*1000)/1000));
        const drawXRatio = getDrawXRatio();
        const drawHeight = 50;
        const drawYRatio = 0.25/drawHeight;
        
        const start = parseTime(e.DestinationArrival.Start)/60/1000;
        const d = [[(start-1-nextBestDeparturesBoundsMinutes[0])/drawXRatio, drawHeight]];
        for (let i=0; i<e.DestinationArrival.Histogram.length; i++) {
            const pos = start+i;
            if (pos < nextBestDeparturesBoundsMinutes[0] || pos > nextBestDeparturesBoundsMinutes[1]) continue;
            d.push([(pos-nextBestDeparturesBoundsMinutes[0])/drawXRatio, drawHeight-e.DestinationArrival.Histogram[i]/drawYRatio]);
        }
        d.push([(start+e.DestinationArrival.Histogram.length-nextBestDeparturesBoundsMinutes[0])/drawXRatio, drawHeight]);
        return 'M '+d.map(p => p.join(' ')).join('L')+' Z';
    }

    function meanDestinationArrival(e: Edge) {
        return simpleTime(parseTime(e.DestinationArrival.Mean));
    }

    function meanDestinationArrivalPos(e: Edge) {
        return (parseTime(e.DestinationArrival.Mean)/60/1000-nextBestDeparturesBoundsMinutes[0])/getDrawXRatio();
    }

    function twoSigmaDestinationArrival(e: Edge) {
        const bounds = getBounds(e);
        return [
            {pos: (bounds[0]-nextBestDeparturesBoundsMinutes[0])/getDrawXRatio(), time: simpleTime(bounds[0]*60*1000)},
            {pos: (bounds[1]-nextBestDeparturesBoundsMinutes[0])/getDrawXRatio(), time: simpleTime(bounds[1]*60*1000)}
        ]
    }

    function isShortestPath(d: Edge) {
        return selection.edge?.ShortestPath.length > 0 && d.ID == selection.edge?.ShortestPath[0].EdgeID;
    }

    function pushHistory(edge: Edge) {
        selectedEdgeHistory = [...selectedEdgeHistory, edge];
        selectEdge(edge.ID);
    }

    function popHistory() {
        if (selectedEdgeHistory.length > 1) {
            selectedEdgeHistory.pop();
            selectEdge(selectedEdgeHistory[selectedEdgeHistory.length-1].ID);
            selectedEdgeHistory = selectedEdgeHistory;
        } else if (selection.edge?.ReverseShortestPath.length > 0) {
            selectEdge(selection.edge.ReverseShortestPath[0].EdgeID);
        } else if (selection.edge?.From.SpaceAxis == data.From.ID) {
            selectStation(data.From.ID);
        }
    }

    function selectedStationName(s: Selection) {
        return selection.edge
            ? stationResolver(selection.edge?.To.SpaceAxis).Name
            : selection.station?.Name;
    }

    function tryDate(old: Date, n: Date, minTime?: Date, maxTime?: Date) {
        const candidate = new Date(old.getFullYear(), old.getMonth(), old.getDate(), n.getUTCHours(), n.getUTCMinutes());
        if (minTime && maxTime && (candidate < minTime || candidate > maxTime)) {
            return undefined;
        }
        return candidate;
    }

    function dayOffset(ref: Date, offset: number) {
        return new Date(ref.getTime()+offset*1000*60*60*24);
    }

    function updateTime(e: any) {
        const old: Date = selection.from;
        const n: Date = e.target.valueAsDate;
        if (!n) return;
        const minTime = new Date(data.MinTime);
        const maxTime = new Date(data.MaxTime);
        let candidate = tryDate(old, n, minTime, maxTime);
        if (!candidate) {
            candidate = tryDate(dayOffset(old, 1), n, minTime, maxTime);
        }
        if (!candidate) {
            candidate = tryDate(dayOffset(old, -1), n, minTime, maxTime);
        }
        if (!candidate) {
            candidate = tryDate(old, n);
        }
        selection.from = candidate;
    }
</script>

<div id="details"><div>
    <div class="refresh"><a href="javascript:void(0)" on:click={doRefresh}>
        <span class="indicator {loading ? 'loading' :''}"><span class="micon">autorenew</span></span>
    </a></div>
    {#if error}
        <p class="error">{$t('c.error')}: {$t('c.'+error)}</p>
    {/if}

    {#if selection.edge}
        <div class="train">
            <h4>
                <span class="label">{@html label(selection.edge, true)}</span>
                <span class="destination">
                    {#if selection.edge.Line && selection.edge.Line.Direction}
                    <span class="micon">east</span> {selection.edge.Line.Direction}
                    {/if}
                </span>
            </h4>
        </div>
    
        <div class="arrdep">
            <span class="left">
                {#if selectedEdgeHistory.length > 1 || selection.edge?.ReverseShortestPath.length > 0 || selection.edge?.From.SpaceAxis == data.From.ID}
                    <a href="javascript:void(0)" on:click={() => popHistory()} class="back"><span class="micon">arrow_back_ios_new</span></a>
                {/if}
            </span>
            <span class="dep">
                {@html departure(selection.edge)}
                {#if selection.edge.Cancelled}<span class="cancelled">({$t('c.cancelled')})</span>{/if}
                <br />{stationResolver(selection.edge.From.SpaceAxis).Name}
            </span>
            <svg viewBox="0 0 50 10" class="miniature">
                <path d="M 10,5 L40,5" class="edge type-{type(selection.edge)} redundant-false"/>
            </svg>
            <span class="arr">
                {@html arrival(selection.edge)}
                {#if selection.edge.Cancelled}<span class="cancelled">({$t('c.cancelled')})</span>{/if}
                <br />{stationResolver(selection.edge.To.SpaceAxis).Name}
            </span>
            <span class="right"></span>
        </div>
        {#if selection.edge.Message}
            <div class="message">
                <span class="micon">info</span> {selection.edge.Message}
            </div>
        {/if}
    {/if}

    {#if nextBestDepartures}
        <h4>{$t('c.next_best_departures')} {selectedStationName(selection)}, <input type="time" value={simpleTime(selection.from)} on:change={updateTime}></h4>
        <table class="next-best-departures">
            <tr>
                <th>{$t('c.header_departure')}</th>
                <th>{$t('c.header_direction')}</th>
                <th>{$t('c.header_destination_arrival')}</th>
            </tr>
        {#each nextBestDepartures as d}
            <tr class="{isShortestPath(d) ? 'shortest' : (d.Redundant ? 'redundant' : '')}" on:click={() => pushHistory(d)}>
                <td>{@html departure(d)}</td>
                <td class="forcewrap">
                    {d.Line?.ID == selection.edge?.Line?.ID ? $t('c.stay_on') : ''}
                    <span>{@html label(d, true)}</span>
                    <span>
                        {#if d.Line?.Direction}
                        <span class="micon">east</span> {d.Line.Direction}
                        {/if}
                        {#if d.Cancelled}<span class="cancelled">({$t('c.cancelled')})</span>{/if}
                    </span>
                    {#if d.Line?.Type == 'Foot' && d.ShortestPath.length > 0}
                    <span>– {@html label(edgeResolver(d.ShortestPath[0].EdgeID), true)}</span>
                    {/if}
                </td>
                <td class="nowrap">
                    {#if hasDistribution(d)}
                        <svg width="100" height="70" class="histogram-canvas">
                            <path d={histogram(d)} class="histogram" />
                            {#each [twoSigmaDestinationArrival(d)] as twoSigma}
                                <path d={'M '+twoSigma[0].pos+' 50 v 5'} class="histogram-pointer twosigma" />
                                <text x={twoSigma[0].pos} y="65" class="histogram-label label twosigma" style="text-anchor:end;">{twoSigma[0].time}</text>
                                <path d={'M '+twoSigma[1].pos+' 50 v 5'} class="histogram-pointer twosigma" />
                                <text x={twoSigma[1].pos} y="65" class="histogram-label label twosigma" style="text-anchor:start;">{twoSigma[1].time}</text>
                                <path d={'M '+twoSigma[0].pos+' 52.5 H '+twoSigma[1].pos+''} class="histogram-pointer twosigma" />
                            {/each}
                            <path d={'M '+meanDestinationArrivalPos(d)+' 50 v 5'} class="histogram-pointer" />
                            <text x={meanDestinationArrivalPos(d)} y="65" class="histogram-label label">{meanDestinationArrival(d)}</text>
                        </svg>
                    {:else}
                        <span class="micon">flag</span>
                        {simpleTime(d.EarliestDestinationArrival)}
                        {#if isShortestPath(d)}<span class="micon">speed</span>{/if}
                    {/if}
                </td>
            </tr>    
        {/each}
        </table>
        {#if nextBestDepartures.length == 0}
            <p>{$t('c.no_known_connections')}</p>
        {/if}
    {/if}

    <a href="?{optionsQueryString(store)}&form" class="submit">
        {$t('c.modify_query')}
    </a>

    <div class="legend">
        <h4>{$t('c.legend')}</h4>
        <svg viewBox="125 1430 500 160" style="width: 100%">
            <g id="legend">
                <g>
                    <path d="M 150,1460 L350,1460" class="edge selected redundant-false" />
                    <text x="150" y="1450">
                        {$t('c.fastest_route')}
                    </text>
                </g>
                <g transform="translate(250, 0)">
                    <path d="M 150,1460 L350,1460" class="edge provider-shortest-path redundant-false" />
                    <text x="150" y="1450">
                        {$t('c.provider_recommended_route')}
                    </text>
                </g>
                <g transform="translate(0, 40)">
                    <path d="M 150,1460 L350,1460" class="edge redundant-true" />
                    <text x="150" y="1450">
                        {$t('c.redundant_connection')}
                    </text>
                </g>
                <g transform="translate(250, 40)">
                    <path d="M 150,1460 L350,1460" class="edge redundant-true cancelled-true" />
                    <text x="150" y="1450">
                        {$t('c.cancelled_connection')}
                    </text>
                </g>
                <g transform="translate(0, 80)">
                    <path d="M 150,1460 L350,1460" class="edge type-nationalExpress redundant-false" />
                    <text x="150" y="1450">
                        {$t('c.long_distance')}
                    </text>
                </g>
                <g transform="translate(250, 80)">
                    <path d="M 150,1460 L350,1460" class="edge type-regional redundant-false" />
                    <text x="150" y="1450">
                        {$t('c.regional')}
                    </text>
                </g>
                <g transform="translate(0, 120)">
                    <path d="M 150,1460 L350,1460" class="edge type-bus redundant-false" />
                    <text x="150" y="1450">
                        {$t('c.local')}
                    </text>
                </g>
                <g transform="translate(250, 120)">
                    <path d="M 150,1460 L350,1460" class="edge type-Foot redundant-false" />
                    <text x="150" y="1450">
                        {$t('c.pedestrian')}
                    </text>
                </g>
            </g>
        </svg>
    </div>
    <Footer />
</div></div>