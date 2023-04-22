<script lang="ts">
    import { t } from '$lib/translations';
    import { onMount, tick } from "svelte";
    import { setFromApi, store } from "../store"
    import { handleHttpErrors, optionsQueryString } from "../query"
    import {parseTime, simpleTime, label, type, departure, arrival, liveDataDeparture, liveDataArrival} from './labels';
    import type {Response, Edge, Coord, Station} from './types';
    import panzoom from 'panzoom'
    import Details from './details.svelte'

    let loading = true;
    let query = store;
    let data: Response;
    let error: string | undefined;
    let currentSelected: Edge = undefined;
    let currentSelectedShortestPath: Edge[] = [];
    const arrowMargin = 25;

    function fetchTimespace(): void {
        fetch(import.meta.env.VITE_TSTP_API+'timespace?'+optionsQueryString(query))
        .then(handleHttpErrors)
        .then(d => {
            data = d;
            console.log(data);
            setFromApi(data);
            loading = false;
            tick().then(() => {
                if (currentSelected && data.Edges[currentSelected.ID]) {
                    selectEdge(currentSelected.ID);
                } else {
                    selectEdge(data.DefaultShortestPathID);
                }                
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
        if (!edgeId) return;
        setSelectedForDependents(false);
        currentSelected = data.Edges[edgeId];
        console.log('cur', currentSelected, edgeId);
        currentSelectedShortestPath = setSelectedForDependents(true);
        console.log(currentSelectedShortestPath);
    }

    function setActive(selected: boolean): void {
        const e = <SVGPathElement><any>document.getElementById(currentSelected.ID+'-toucharea');
        if (!e) return;
        if (selected) {
            e.className.baseVal += " active";
        } else {
            e.className.baseVal =  e.className.baseVal.replace(" active", "");
        }
    }

    function setSelectedForDependents(selected: boolean): Edge[] {
        if (!currentSelected) {
            return [];
        }
        setActive(selected);
        let all = [];
        let previous = currentSelected;
        let next = currentSelected;
        while(true) {
            setSelectedForEdge(next, selected);
            setSelectedForStationEdge(previous, next, selected);
            all.push(next);
            if (next.ShortestPath.length == 0) break;
            previous = next;
            next = data.Edges[next.ShortestPath[0].EdgeID];
            if (!next) break;
        }
        next = currentSelected;
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
        console.log('selected ', id);
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
        
        let delta = (parseTime(coord.TimeAxis)-parseTime(data.MinTime))/1000;
        return delta/data.TimeAxisDistance*(data.TimeAxisSize-100)+100;
    }

    function stationResolver(id: string): Station {
        return data.Stations[id];
    }

    function stationsList(): Station[] {
        return Object.values(data.Stations);
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
    <span class="indicator {loading ? 'loading' :''}"><span class="micon">autorenew</span></span>
    <p>{$t('c.data_retrieval_waiting')}</p>
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
{#if data}
{#each stationsList() as s (s.ID)}
<text x="{x(s.Coord)}" y="{y(s.Coord)}" class="station-label">
    {s.Name}
    <title>{s.ID}</title>
</text>
{/each}
{#each data.TimeIndicators as t}
<text x="5" y="{y(t)}" class="time-label">
    {simpleTime(t.TimeAxis)}
</text>
{/each}
{#each data.SortedEdges.map(id => data.Edges[id]) as e (e.ID)}
{#if !e.Discarded}
<path id="{e.ID}" d="M {x(e.From)},{y(e.From)} L{x(e.To)},{y(e.To)}"
    class="edge type-{type(e)} redundant-{e.Redundant} cancelled-{e.Cancelled} {e.ProviderShortestPath ? 'provider-shortest-path' : ''}"
    />
<path id="{e.ID}-toucharea" d="M {x(e.From)},{y(e.From)} L{x(e.To)},{y(e.To)}"
    class="edge-toucharea" on:click={selectListener}
    />
{/if}
{/each}
{#each currentSelectedShortestPath as e (e.ID)}
{#if !e.Discarded}
<text id="{e.ID}-label" class="label type-{type(e)} label-{e.ID}">
    <textPath href="#{e.ID}" startOffset="50%">
        {@html label(e, false)}
    </textPath>
</text>
<text x="{x(e.To)}" y="{y(e.To)}" id="{e.ID}-arrival"
    class="arrival type-{type(e)} label-{e.ID} {liveDataArrival(e)}">
    {arrival(e)}
</text>
<text x="{x(e.From)}" y="{y(e.From)}" id="{e.ID}-departure"
    class="departure type-{type(e)} label-{e.ID} {liveDataDeparture(e)}">
    {departure(e)}
</text>
{#if e.PreviousArrival}<text x="{x(e.To)}" y="{y(e.To)-arrowMargin}" class="previous-next-arrow" on:click={() => selectEdge(e.PreviousArrival)}>▲</text>{/if}
{#if e.PreviousDeparture}<text x="{x(e.From)}" y="{y(e.From)-arrowMargin}" class="previous-next-arrow" on:click={() => selectEdge(e.PreviousDeparture)}>▲</text>{/if}
{#if e.NextArrival}<text x="{x(e.To)}" y="{y(e.To)+arrowMargin}" class="previous-next-arrow" on:click={() => selectEdge(e.NextArrival)}>▼</text>{/if}
{#if e.NextDeparture}<text x="{x(e.From)}" y="{y(e.From)+arrowMargin}" class="previous-next-arrow" on:click={() => selectEdge(e.NextDeparture)}>▼</text>{/if}
{/if}
{/each}
{/if}

</svg>
</div>
<Details currentSelected={currentSelected} loading={loading} doRefresh={refresh} stationResolver={stationResolver} error={error}/>