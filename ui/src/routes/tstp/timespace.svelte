<script lang="ts">
    import { page } from '$app/stores';
    import { t } from '$lib/translations';
    import { onMount, tick } from "svelte";
    import { defaultDatetime, setFromApi, store } from "../store"
    import { handleHttpErrors, optionsQueryString } from "../query"
    import {parseTime, simpleTime, label, type, departure, arrival, redundant} from './labels';
    import {type Response, type Edge, type Coord, Selection} from './types';
    import panzoom from 'panzoom'
    import Details from './details.svelte'

    let loading = true;
    let query = store;
    let data: Response;
    let error: string | undefined;
    let displayTsd = false;
    let selection: Selection = new Selection();
    let selectedShortestPath: Edge[] = [];
    const arrowMargin = 25;

    function fetchTimespace(): void {
        fetch(import.meta.env.VITE_TSTP_API+'timespace?'+optionsQueryString(query, defaultDatetime))
        .then(handleHttpErrors)
        .then(d => {
            setSelectedForDependents(false);
            data = d;
            setFromApi(data);
            loading = false;
            tick().then(() => {
                if (selection.edge && data.Edges[selection.edge.ID]) {
                    selectEdge(selection.edge.ID);
                } else if (selection.station && data.Stations[selection.station.ID]) {
                    selectStation(selection.station.ID);
                } else {
                    selectStation(data.From.ID);
                }      
                window.location.hash = 'details';
            });
        })
        .catch((err) => {
            loading = false;
            console.log('Error:', err);
            error = err.message && err.message.startsWith('error_http_') ? err.message : 'error_unknown';
        });
    }

    function refresh() {
        error = undefined;
        loading = true;
        fetchTimespace();
    }

    function selectEdge(edgeId: string | undefined): void {
        setSelectedForDependents(false);
        if (!edgeId) {
            selectedShortestPath = [];
            return;
        }
        selection = Selection.fromEdge(data.Edges[edgeId]);
        console.log('cur', selection.edge, edgeId);
        selectedShortestPath = setSelectedForDependents(true);
    }

    function selectStation(stationId: string) {
        if (!stationId) return;
        selectEdge(undefined);
        selection = Selection.fromStation(data.Stations[stationId]);
    }

    function setActive(selected: boolean): void {
        const e = <SVGPathElement><any>document.getElementById(selection.edge.ID+'-toucharea');
        if (!e) return;
        if (selected) {
            e.className.baseVal += " active";
        } else {
            e.className.baseVal =  e.className.baseVal.replace(" active", "");
        }
    }

    function setSelectedForDependents(selected: boolean): Edge[] {
        if (!selection.edge) {
            return [];
        }
        setActive(selected);
        let all = [];
        let previous = selection.edge;
        let next = selection.edge;
        while(true) {
            setSelectedForEdge(next, selected);
            setSelectedForStationEdge(previous, next, selected);
            all.push(next);
            if (next.ShortestPath.length == 0) break;
            previous = next;
            next = data.Edges[next.ShortestPath[0].EdgeID];
            if (!next) break;
        }
        next = selection.edge;
        while(next.ReverseShortestPath.length > 0) {
            previous = next;
            next = data.Edges[next.ReverseShortestPath[0].EdgeID];
            if (!next) break;
            setSelectedForEdge(next, selected);
            setSelectedForStationEdge(previous, next, selected);
            all.push(next);
        }
        return all;
    }

    function setSelectedForEdgeId(edgeId: string, selected: boolean): void {
        const e = <SVGPathElement><any>document.getElementById(edgeId);
        if (!e) return;
        if (selected) {
            e.className.baseVal += " selected";
        } else {
            e.className.baseVal =  e.className.baseVal.replace(" selected", "").replace(" active", "");
        }    
    }

    function setSelectedForEdge(edge: Edge, selected: boolean): void {
        if (!edge || edge.Discarded) return;
        setSelectedForEdgeId(edge.ID, selected);
    }
    
    
    function setSelectedForStationEdge(previous: Edge, next: Edge, selected: boolean): void {
        if (!previous || !next || previous.Discarded || next.Discarded) return;
        const edgeId = previous.ID+'_'+next.ID+'_station';
        if (document.getElementById(edgeId)) {
            setSelectedForEdgeId(edgeId, selected);
        }
    }

    function selectListener(evt: any): void {
        const id = evt.target.id.replace('-toucharea', '');
        selectEdge(id);
    }

    function x(coord: Coord): number {
        if (coord.SpaceAxis == "") {
            return 0;
        }
        return data.Stations[coord.SpaceAxis].SpaceAxis/data.MaxSpace*(data.SpaceAxisSize-100)+50;
    }

    function y(coord: Coord): number {
        if (parseTime(coord.TimeAxis) == 0) {
            return 50 + data.Stations[coord.SpaceAxis].SpaceAxisHeap*20;
        }
        return yByTs(parseTime(coord.TimeAxis));       
    }

    function yByTs(unixTs: number) {
        let delta = (unixTs-parseTime(data.MinTime))/1000;
        return delta/data.TimeAxisDistance*(data.TimeAxisSize-100)+100;
    }

    function randomTipId() {
        const tipCount = 4;
        return 'tip_'+Math.floor(Math.random()*tipCount);
    }

    function showTsd() {
        displayTsd = true;
    }

    onMount(() => {
        fetchTimespace();
        const e = document.getElementById('timespace-canvas')
        if (e) panzoom(e, {
                maxZoom: 7,
                minZoom: 1,
                bounds: true,
                boundsPadding: 1,
                zoomDoubleClickSpeed: 3,
                onTouch: function(e) {
                    if ((<any>e.target)?.id == 'timespace-canvas') return true; 
                }
        });
    })
</script>

<div class="loading-screen" style="display: {!data ? 'block' : 'none'};">
    <span class="refresh"><a href="javascript:void(0)" on:click={refresh}>
        <span class="indicator {loading ? 'loading' :''}"><span class="micon">autorenew</span></span>
    </a></span>
    <p>{$t('c.data_retrieval_waiting')}</p>
    <p><span class="micon">tips_and_updates</span> {$t('c.tip')+': '+$t('c.'+randomTipId())}</p>
</div>
<div id="timespace-container">

<svg id="timespace-canvas" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="xMidYMin"
    viewBox="0 0 1500 1500">
<defs>
    <marker id="dot-redundant-false" viewBox="0 0 10 10" refX="5" refY="5"
        markerWidth="4" markerHeight="4" filter="none">
        <circle cx="5" cy="5" r="3" />
    </marker>
    <marker id="dot-redundant-true" viewBox="0 0 10 10" refX="5" refY="5"
        markerWidth="4" markerHeight="4" filter="none">
        <circle cx="5" cy="5" r="3" />
    </marker>
    <filter id="erode">
        <feMorphology operator="erode" radius="8"/>
    </filter>
</defs>
{#if data && displayTsd}
{#each Object.values(data.Stations) as s (s.ID)}
<text x="{x(s.Coord)}" y="{y(s.Coord)}" class="station-label" on:click={() => selectStation(s.GroupID || s.ID)}>
    {s.Name}
    <title>{s.ID}</title>
</text>
{/each}
{#each data.TimeIndicators as t}
<text x="5" y="{y(t)}" class="time-label">
    {simpleTime(t.TimeAxis)}
</text>
{/each}
{#if selection.station && selection.from}
<path id="station-toucharea" d="M {x(selection.station.Coord)},{yByTs(selection.from?.getTime())} L{x(selection.station.Coord)},{yByTs(selection.from?.getTime()+3600*1000)}"
    class="edge-toucharea active" />
{/if}
{#each data.SortedEdges.map(id => data.Edges[id]) as e (e.ID)}
{#if !e.Discarded}
<path id="{e.ID}" d="M {x(e.From)},{y(e.From)} L{x(e.To)},{y(e.To)}"
    class="edge type-{type(e)} redundant-{redundant(e)} cancelled-{e.Cancelled} {e.ProviderShortestPath ? 'provider-shortest-path' : ''}"
    />
<path id="{e.ID}-toucharea" d="M {x(e.From)},{y(e.From)} L{x(e.To)},{y(e.To)}"
    class="edge-toucharea" on:click={selectListener}
    />
{/if}
{/each}
{#each selectedShortestPath as e (e.ID)}
{#if !e.Discarded}
<text id="{e.ID}-label" class="label type-{type(e)} label-{e.ID}">
    <textPath href="#{e.ID}" startOffset="50%">
        {@html label(e, false)}
    </textPath>
</text>
<text x="{x(e.To)}" y="{y(e.To)}" id="{e.ID}-arrival"
    class="arrival type-{type(e)} label-{e.ID}">
    {@html arrival(e, 'tspan')}
</text>
<text x="{x(e.From)}" y="{y(e.From)}" id="{e.ID}-departure"
    class="departure type-{type(e)} label-{e.ID}">
    {@html departure(e, 'tspan')}
</text>
{#if e.PreviousArrival}<text x="{x(e.To)}" y="{y(e.To)-arrowMargin}" class="previous-next-arrow" on:click={() => selectEdge(e.PreviousArrival)}>▲</text>{/if}
{#if e.PreviousDeparture}<text x="{x(e.From)}" y="{y(e.From)-arrowMargin}" class="previous-next-arrow" on:click={() => selectEdge(e.PreviousDeparture)}>▲</text>{/if}
{#if e.NextArrival}<text x="{x(e.To)}" y="{y(e.To)+arrowMargin}" class="previous-next-arrow" on:click={() => selectEdge(e.NextArrival)}>▼</text>{/if}
{#if e.NextDeparture}<text x="{x(e.From)}" y="{y(e.From)+arrowMargin}" class="previous-next-arrow" on:click={() => selectEdge(e.NextDeparture)}>▼</text>{/if}
{/if}
{/each}
{/if}

</svg>

{#if data && !displayTsd}
    <div id="show-tsd">
    <a href="javascript:void(0)" on:click={showTsd} class="submit">
        {$t('c.show_tsd')}
    </a>
    </div>
{/if}
</div>
<Details bind:selection={selection} loading={loading} tsdShown={displayTsd} doRefresh={refresh} selectEdge={selectEdge} selectStation={selectStation} data={data} error={error}/>