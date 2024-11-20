<script lang="ts">
    import { page } from '$app/stores';
    import { t } from '$lib/translations';
    import { store } from "../store"
    import { optionsQueryString } from "../query"
    import Footer from '../footer.svelte'; 
    import {label, type, departure, arrival, parseTime, simpleTime, redundant} from './labels';
    import type { Edge, Station, Response } from './types';
    import { calcNextDepartureIndex, getStationsInGroup, hasDistribution, walkingDistance, walkingDurationMs } from './shortestPaths';

    export let selection: Edge;
    export let loading: boolean;
    export let tsdShown: boolean;
    export let doRefresh: () => void;
    export let selectEdge: (id: string) => void;
    export let selectStation: (id: string) => void;
    export let data: Response;
    export let error: string | undefined;

    let selectedEdgeHistory: Edge[] = [];

    let nextBestDepartures: Edge[] | undefined = undefined;
    let nextBestDeparture: Edge | undefined = undefined;
    let nextBestDeparturesBoundsMinutes: number[] = [0,0];
    const defaultNextBestDepartures = 5;
    const defaultNegativeTransferMinutes = 5;
    let numDepartures = defaultNextBestDepartures;
    let negativeTransferMinutes = defaultNegativeTransferMinutes;

    $: {
        console.log("upd");
        numDepartures = defaultNextBestDepartures;
        negativeTransferMinutes = defaultNegativeTransferMinutes;
        nextBestDepartureForSelection(selection);
    }

    function displayEarlierDepartures() {
        negativeTransferMinutes += 30;
        nextBestDepartureForSelection(selection);
    }

    function displayLaterDepartures() {
        numDepartures += defaultNextBestDepartures;
        nextBestDepartureForSelection(selection);
    }

    function nextBestDepartureForSelection(selection: Edge) {
        if (selection.edge) {
            console.log(selection.edge);
            nextBestDeparturesForEdge();
            rectifyEdgeHistory();
        } else if (selection.station) {
            console.log(selection.station);
            nextBestDeparturesForStation();
        }
    }

    function nextBestDeparturesForEdge() {
        const station = stationResolver(selection.edge.To.SpaceAxis);
        if (selection.edge.To.SpaceAxis == data.To.ID || station.GroupID == stationResolver(data.To.ID).GroupID) {
            nextBestDepartures = undefined;
            return;
        }
        if (!selection.from) {
            selection.from = new Date(parseTime(selection.edge.Actual.Arrival));
        }
        updateNextBestDepartures(station, selection.from);
    }

    function nextBestDeparturesForStation() {
        if (!selection.from) {
            let n = new Date().getTime();
            if (selectedEdgeHistory.length > 0) {
                n = parseTime(selectedEdgeHistory[selectedEdgeHistory.length-1].Actual.Departure);
            } else if (n < new Date(data.MinTime).getTime() || n > new Date(data.MaxTime).getTime()) {
                n = new Date(data.MinTime).getTime()+negativeTransferMinutes*60*1000;
            }            
            selection.from = new Date(n);
        }
        selectedEdgeHistory = [];
        updateNextBestDepartures(selection.station, selection.from);
    }

    function getBounds(e: Edge): number[] {
        const h = e.DestinationArrival.Histogram;
        if (!h) {
            const fallback = parseTime(e.EarliestDestinationArrival)/60/1000;
            return [fallback, fallback];
        }
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
        if (lower != undefined && upper != undefined) {
            return [lower, upper];
        }
        return [start, start];
    }

    function updateNextBestDepartures(station: Station, time: Date) {      
        const candidates = [];
        
        let lowerBound = undefined;
        let upperBound = undefined;
        if (selection.edge) {
            let b = getBounds(selection.edge);
            lowerBound = b[0];
            upperBound = b[1];
        }
        let shortestPathFound = false;
        let relevantStations = getStationsInGroup(data, station);
        let indices = new Array(relevantStations.length).fill(0);
        while (true) {
            if (candidates.length >= numDepartures && shortestPathFound || candidates.length >= numDepartures*numDepartures) {
                break;
            }
            let nextDepartureIndex = calcNextDepartureIndex(station, relevantStations, indices, (_) => time.getTime()-negativeTransferMinutes*60*1000, edgeResolver, stationResolver);
            if (nextDepartureIndex == undefined) {
                break;
            }
            const e = edgeResolver(relevantStations[nextDepartureIndex].BestDepartures[indices[nextDepartureIndex]]);
            let departure = parseTime(e.Actual.Departure);
            candidates.push(e);
            indices[nextDepartureIndex]++;
            if (!hasDistribution(e)) continue;
            if (departure >= time.getTime()+walkingDurationMs(station.ID, e.From.SpaceAxis, stationResolver) && (!shortestPathFound || nextBestDeparture.DestinationArrival?.Mean > e.DestinationArrival?.Mean)) {
                shortestPathFound = true;
                nextBestDeparture = e;
            }
            const bounds = getBounds(e);
            if (lowerBound == undefined || bounds[0] < lowerBound) {
                lowerBound = bounds[0];
            }
            if (upperBound == undefined || bounds[1] > upperBound) {
                upperBound = bounds[1];
            }
        }
        candidates.sort((a, b) => {
            return parseTime(a.Actual.Departure)-parseTime(b.Actual.Departure);
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

    function getDrawXRatio() {
        const drawWidth = 100;
        return (nextBestDeparturesBoundsMinutes[1]-nextBestDeparturesBoundsMinutes[0])/drawWidth;
    }

    function histogram(e: Edge) {
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

    function getMode(e: Edge) {
        const h = e.DestinationArrival.Histogram;
        const start = parseTime(e.DestinationArrival.Start)/60/1000;
        let modeY = undefined;
        let modeX = 0;
        for (let i=0;i<h.length;i++) {
            if (modeY == undefined || modeY < h[i]) {
                modeY = h[i];
                modeX = i;
            }
        }
        return start+modeX;
    }

    function toPosAndTime(minutes: number) {
        return {pos: (minutes-nextBestDeparturesBoundsMinutes[0])/getDrawXRatio(), time: simpleTime(minutes*60*1000)}
    }

    function destinationArrivalMetrics(e: Edge) {
        const bounds = getBounds(e);
        return {
            "p5": toPosAndTime(bounds[0]),
            "p95": toPosAndTime(bounds[1]),
            "mode": toPosAndTime(getMode(e)),
            "mean": toPosAndTime(parseTime(e.DestinationArrival.Mean)/60/1000)
        }
    }

    function isShortestPath(d: Edge) {
        return d == nextBestDeparture;
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
        } else {
            selectStation(selection.edge?.From.SpaceAxis);
            selectedEdgeHistory.pop();
        }
    }

    function selectedStationName(s: Selection) {
        return selection.edge
            ? stationResolver(selection.edge?.To.SpaceAxis).Name
            : selection.station?.Name;
    }

    function selectedStationId(s: Selection) {
        return selection.edge?.To.SpaceAxis || selection.station?.ID
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

    function walkingDistanceRounded(fromId: string, toId: string): number {
        return Math.max(Math.round(walkingDistance(fromId, toId, stationResolver)/100)*100, 100);
    }

    function shareDestArr() {
        const canvas = document.createElement("canvas");
        canvas.width = 1100;
        canvas.height = 180;
        const ctx = canvas.getContext("2d");
        const svgEl = document.getElementById('destarr-histogram');
        if (svgEl && ctx) {
            var sheets = document.styleSheets;
            var styleStr = '';
            Array.prototype.forEach.call(sheets, function(sheet){
                try { 
                    styleStr += Array.prototype.reduce.call(sheet.cssRules, function(a, b) {
                            return a + b.cssText;
                        }, "");
                }
                catch (e) {
                    console.log(e);
                }
            });
            document.getElementById('inherit-style')!.innerHTML = styleStr;
            const svgString = new XMLSerializer().serializeToString(svgEl);
            const domUrl = self.URL || self.webkitURL || self;
            const img = new Image();
            const svg = new Blob([svgString], {type: "image/svg+xml;charset=utf-8"});
            const url = domUrl.createObjectURL(svg);
            img.onload = function() {
                ctx.drawImage(img, 0, 0);
                canvas.toBlob(async (blob) => {
                    if (!blob) {
                        return;
                    }
                    const files = [new File([blob], 'destination_arrival.png', { type: blob.type })];
                    const metrics = destinationArrivalMetrics(selection.edge);
                    const destName = stationResolver(data.To.ID).Name;
                    const shareData = {
                        text: t.get('c.explain_share_mean', {destination: destName, time: metrics.mean.time}) + ' ' +
                            t.get('c.explain_share_mode', {time: metrics.mode.time}) + ' ' +
                            t.get('c.explain_share_95th', {time: metrics.p95.time}),
                        title: t.get('c.explain_share_title', {destination: destName}),
                        files,
                    };
                    if (navigator.canShare(shareData)) {
                        try {
                            await navigator.share(shareData);
                        } catch (err) {
                            console.error(err);
                        }
                    } else {
                        console.warn('Sharing not supported', shareData);
                    }
                });
                domUrl.revokeObjectURL(url);
            }
            img.src = url;
        }
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
                {#if selection.edge.Line && selection.edge.Line.Type.includes('national')}
                <a
                    target="_blank"
                    class="unstyled-link blue-link"
                    href="https://bahnapp.online/journey/coaches/?trainId={selection.edge.Line.Name.replace(' ', '+')}&stationId={selection.edge.From.SpaceAxis}&departureTime={new Date(parseTime(selection.edge.Actual.Departure)).toISOString()}&initialDepartureTime={new Date(parseTime(selection.edge.Planned.Departure)).toISOString()}"
                    title="{$t('c.coach_sequence')}" aria-label="{$t('c.coach_sequence')}">
                        <span class="micon" aria-label="{$t('c.coach_sequence')}">airline_seat_recline_normal</span>
                </a>
                {/if}
            </h4>
        </div>
    
        <div class="arrdep">
            <span class="left">
                {#if selectedEdgeHistory.length > 0}
                    <a href="javascript:void(0)" on:click={() => popHistory()} class="back"><span class="micon">arrow_back_ios_new</span></a>
                {/if}
            </span>
            <span class="dep" title="{stationResolver(selection.edge.From.SpaceAxis).ID}">
                {@html departure(selection.edge)}
                {#if selection.edge.Cancelled}<span class="cancelled">({$t('c.cancelled')})</span>{/if}
                <br />{stationResolver(selection.edge.From.SpaceAxis).Name}
            </span>
            <svg viewBox="0 0 50 10" class="miniature">
                <path d="M 10,5 L40,5" class="edge type-{type(selection.edge)} redundant-false"/>
            </svg>
            <span class="arr" title="{stationResolver(selection.edge.To.SpaceAxis).ID}">
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
        {#if hasDistribution(selection.edge)}
        {#each [destinationArrivalMetrics(selection.edge)] as metrics}
            <a href="javascript:void(0)" on:click={() => shareDestArr()} class="unstyled-link destarr" style="position:relative;display:block">
                <svg width="1100" height="180" class="histogram-canvas" id="destarr-histogram" viewBox="-225 0 550 90" preserveAspectRatio="xMidYMid slice" style="max-width:100%;max-height:90px;">
                    <defs><style id="inherit-style"></style></defs>
                    <rect x="-225" y="0" width="550" height="90" fill="white" />
                    <path d={histogram(selection.edge)} class="histogram" />
                    <path d={'M '+metrics.p5.pos+' 50 v 5'} class="histogram-pointer twosigma" />
                    <path d={'M '+metrics.p95.pos+' 50 v 5'} class="histogram-pointer twosigma" />
                    <text x={metrics.p95.pos} y="67" class="histogram-label label twosigma" style="text-anchor:start;">{metrics.p95.time} <tspan class="explain">&nbsp;({$t('c.explain_95th', {time: metrics.p95.time})})</tspan></text>
                    <path d={'M '+metrics.p5.pos+' 52.5 H '+metrics.p95.pos+''} class="histogram-pointer twosigma" />
                    <path d={'M '+metrics.mode.pos+' 19 v 5'} class="histogram-pointer twosigma" />
                    <text x={metrics.mode.pos} y="7" class="histogram-label label twosigma" style="dominant-baseline:hanging;">{metrics.mode.time}</text>
                    <text x={metrics.mode.pos+20} y="7" class="histogram-label label twosigma explain" style="dominant-baseline:hanging;text-anchor:start;">({$t('c.explain_mode', {destination: stationResolver(data.To.ID).Name, time: metrics.mode.time})})</text>
                    <path d={'M '+metrics.mean.pos+' 50 v 5'} class="histogram-pointer" />
                    <text x={metrics.mean.pos} y="67" class="histogram-label label">{metrics.mean.time}</text>
                    <text x={metrics.mean.pos} y="82" class="histogram-label label explain">({$t('c.explain_mean', {destination: stationResolver(data.To.ID).Name, time: metrics.mean.time})})</text>
                </svg>
                <span class="histogram-label explain" style="position: absolute;bottom:50%;right:0">
                    <span class="micon">share</span>
                </span>
            </a>
        {/each}
        {/if}
    {/if}

    {#if nextBestDepartures}
        <h3><span class="small">{$t('c.next_best_departures')}</span><br /><span class="highlight">{selectedStationName(selection)}</span>, <input type="time" value={simpleTime(selection.from)} on:change={updateTime}></h3>
        <table class="next-best-departures">
            <tr><td class="nosep right" colspan="3"><a href="javascript:void(0)" on:click={displayEarlierDepartures} class="submit small">
                {$t('c.earlier_departures')}
            </a></td></tr>
            <tr>
                <th>{$t('c.header_departure')}</th>
                <th>{$t('c.header_direction')}</th>
                <th>{$t('c.header_destination_arrival')}</th>
            </tr>
        {#if selection.from.getTime()-negativeTransferMinutes*60*1000 < new Date(data.MinTime).getTime()}
            <tr><td colspan="3" class="no-conns">{$t('c.no_earlier_connections')}</td></tr>
        {/if}
        {#each nextBestDepartures as d}
            <tr class="{isShortestPath(d) ? 'shortest' : (redundant(d) ? 'redundant' : '')}" on:click={() => pushHistory(d)}>
                <td>{@html departure(d)}</td>
                <td class="forcewrap">
                    <span class="stay-on">{d.Line?.ID == selection.edge?.Line?.ID && selection.edge?.To.SpaceAxis == d.From.SpaceAxis ? $t('c.stay_on') : ''}</span>
                    <span>{@html label(d, true)}</span>
                    <span>
                        {#if d.Line?.Direction}
                        <span class="micon">east</span> {d.Line.Direction}
                        {/if}
                        {#if d.Cancelled}<span class="cancelled">({$t('c.cancelled')})</span>{/if}
                        {#if d.Line?.Type == 'Foot'}{$t('c.walking_to')} {stationResolver(d.To.SpaceAxis).Name}{/if}
                    </span>
                    {#if d.From.SpaceAxis != selectedStationId(selection)}
                    <span class="walking-from">{$t('c.walking_from')} {stationResolver(d.From.SpaceAxis).Name.split(',')[0]}, <span class="micon">directions_walk</span>{walkingDistanceRounded(selectedStationId(selection), d.From.SpaceAxis)}m</span>
                    {/if}
                </td>
                <td class="nowrap">
                    {#if hasDistribution(d)}
                        <svg width="100" height="70" class="histogram-canvas">
                            <path d={histogram(d)} class="histogram" />
                            {#each [destinationArrivalMetrics(d)] as metrics}
                            <path d={'M '+metrics.p5.pos+' 50 v 5'} class="histogram-pointer twosigma" />
                            <!--<text x={metrics.p5.pos} y="67" class="histogram-label label twosigma" style="text-anchor:end;">{metrics.p5.time}</text>-->
                            <path d={'M '+metrics.p95.pos+' 50 v 5'} class="histogram-pointer twosigma" />
                            <text x={metrics.p95.pos} y="67" class="histogram-label label twosigma" style="text-anchor:start;">{metrics.p95.time}</text>
                            <path d={'M '+metrics.p5.pos+' 52.5 H '+metrics.p95.pos+''} class="histogram-pointer twosigma" />
                            <path d={'M '+metrics.mode.pos+' 19 v 5'} class="histogram-pointer twosigma" />
                            <text x={metrics.mode.pos} y="7" class="histogram-label label twosigma" style="dominant-baseline:hanging;">{metrics.mode.time}</text>
                            <path d={'M '+metrics.mean.pos+' 50 v 5'} class="histogram-pointer" />
                            <text x={metrics.mean.pos} y="67" class="histogram-label label">{metrics.mean.time}</text>
                            {/each}
                        </svg>
                    {:else}
                        <span class="micon">flag</span>
                        {simpleTime(d.EarliestDestinationArrival)}
                        {#if isShortestPath(d)}<span class="micon">speed</span>{/if}
                    {/if}
                </td>
            </tr>    
        {/each}
        {#if nextBestDepartures.length < numDepartures}
            <tr><td colspan="3" class="no-conns">{$t('c.no_later_connections')}</td></tr>
        {/if}
        <tr><td class="right" colspan="3"><a href="javascript:void(0)" on:click={displayLaterDepartures} class="submit small">
            {$t('c.later_departures')}
        </a></td></tr>
        </table>
        
    {/if}

    <a href="?{optionsQueryString(store)}&form" class="submit">
        {$t('c.modify_query')}
    </a>
    {#if $page.url.searchParams.get('tsd') != 'no' && tsdShown}
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
    {/if}
    <Footer />
</div></div>