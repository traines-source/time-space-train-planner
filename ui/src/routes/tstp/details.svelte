<script lang="ts">
    import { t } from '$lib/translations';
    import { store } from "../store"
    import { optionsQueryString } from "../query"
    import Footer from '../footer.svelte'; 
    import {label, type, departure, arrival, liveDataDeparture, liveDataArrival, parseTime, simpleTime} from './labels';
    import type { Edge, Station, Response } from './types';

    export let selection: Edge;
    export let loading: boolean;
    export let doRefresh: () => void;
    export let selectEdge: (id: string) => void;
    export let selectStation: (id: string) => void;
    export let data: Response;
    export let error: string | undefined;

    let selectedEdgeHistory: Edge[] = [];

    let nextBestDepartures: Edge[] = [];
    const maxNextBestDepartures = 5;
    const negativeTransferMinutes = 5;
    $: {
        if (selection.edge) {
            updateNextBestDepartures();
            rectifyEdgeHistory();
        }
    }

    function updateNextBestDepartures() {
        if (selection.edge.To.SpaceAxis == data.To.ID) {
            nextBestDepartures = [];
            return;
        }
        const station = stationResolver(selection.edge.To.SpaceAxis);
        const candidates = [];
        const currentTime = parseTime(selection.edge.Actual.Arrival)-negativeTransferMinutes*60*1000;
        for (let i=0; i<station.BestDepartures.length; i++) {
            const e = edgeResolver(station.BestDepartures[i]);
            if (parseTime(e.Actual.Departure) < currentTime) {
                continue;
            }
            candidates.push(e);
            if (candidates.length >= maxNextBestDepartures) {
                break;
            }
        }
        candidates.sort((a, b) => {
            return parseTime(a.Planned.Departure)-parseTime(b.Planned.Departure);
        });
        nextBestDepartures = candidates;
    }

    function rectifyEdgeHistory() {
        if (selectedEdgeHistory[selectedEdgeHistory.length-1] != selection.edge) {
            selectedEdgeHistory = [selection.edge];
            console.log("history reset");
        }
    }

    function edgeResolver(id: string): Edge {
        return data.Edges[id];
    }

    function stationResolver(id: string): Station {
        return data.Stations[id];
    }

    function isShortestPath(d: Edge) {
        return selection.edge.ShortestPath.length > 0 && d.ID == selection.edge.ShortestPath[0].EdgeID;
    }

    function pushEdge(edge: Edge) {
        selectedEdgeHistory = [...selectedEdgeHistory, edge];
        selectEdge(edge.ID);
    }

    function popEdge() {
        if (selectedEdgeHistory.length > 1) {
            selectedEdgeHistory.pop();
            selectEdge(selectedEdgeHistory[selectedEdgeHistory.length-1].ID);
            selectedEdgeHistory = selectedEdgeHistory;
        } else if (selection.edge.ReverseShortestPath.length > 0) {
            selectEdge(selection.edge.ReverseShortestPath[0].EdgeID);
        }
    }
</script>

<div id="details"><div>
    {#if error}
    <p class="error">{$t('c.error')}: {$t('c.'+error)}</p>
    {/if}
    {#if selection.edge}
    <div class="refresh"><a href="javascript:void(0)" on:click={doRefresh}>
        <span class="indicator {loading ? 'loading' :''}"><span class="micon">autorenew</span></span>
    </a></div>
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
            {#if selectedEdgeHistory.length > 1 || selection.edge.ReverseShortestPath.length > 0}
                <a href="javascript:void(0)" on:click={() => popEdge()} class="back"><span class="micon">arrow_back_ios_new</span></a>
            {/if}
        </span>
        <span class="dep">
            <span class="{liveDataDeparture(selection.edge)}">{departure(selection.edge)}</span>
            {#if selection.edge.Cancelled}<span class="cancelled">({$t('c.cancelled')})</span>{/if}
            <br />{stationResolver(selection.edge.From.SpaceAxis).Name}
        </span>
        <svg viewBox="0 0 50 10" class="miniature">
            <path d="M 10,5 L40,5" class="edge type-{type(selection.edge)} redundant-false"/>
        </svg>
        <span class="arr">
            <span class="{liveDataArrival(selection.edge)}">{arrival(selection.edge)}</span>
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

    {#if nextBestDepartures.length > 0}
        <h4>{$t('c.next_best_departures')} {stationResolver(selection.edge.To.SpaceAxis).Name}, ~{simpleTime(selection.edge.Actual.Arrival)}</h4>
    {/if}
    <table class="next-best-departures">
    {#each nextBestDepartures as d}
        <tr class="{isShortestPath(d) ? 'shortest' : (d.Redundant ? 'redundant' : '')}" on:click={() => pushEdge(d)}>
            <td class="nowrap"><span class="{liveDataDeparture(d)}">{departure(d)}</span></td>
            <td class="forcewrap">
                {d.Line && selection.edge.Line && d.Line.ID == selection.edge.Line.ID ? $t('c.stay_on') : ''}
                <span>{@html label(d, true)}</span>
                <span>
                    {#if d.Line && d.Line.Direction}
                    <span class="micon">east</span> {d.Line.Direction}
                    {/if}
                    {#if d.Cancelled}<span class="cancelled">({$t('c.cancelled')})</span>{/if}
                </span>
            </td>
            <td class="nowrap">
                <span class="micon">flag</span>
                {simpleTime(d.EarliestDestinationArrival)}
                {#if isShortestPath(d)}<span class="micon">speed</span>{/if}
            </td>
        </tr>    
    {/each}
    </table>
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