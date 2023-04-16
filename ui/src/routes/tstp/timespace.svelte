<script>
    import { onMount } from "svelte";
    import { setFromApi, Station, store } from "../store"
    import { optionsQueryString } from "../url"
    
    let loading = true;
    let query = store;
    let data;

    function fetchTimespace() {
        fetch(import.meta.env.VITE_TSTP_API+'timespace?'+optionsQueryString(query))
        .then(response => response.json())
        .then(d => {
            data = d;
            console.log(data);
            setFromApi(data);
            console.log(store);
            loading = false;
        })
        .catch((error) => {
            alert('Failed request. Possibly too many requests. Try again later.');
            loading = false;
            console.log(error);
        });
    }

    function parseTime(t) {
        return Math.max(0, Date.parse(t));
    }

    function x(coord) {
        if (coord.SpaceAxis == "") {
            return 0;
        }
	    return data.Stations[coord.SpaceAxis].SpaceAxis/data.MaxSpace*(data.SpaceAxisSize-100)+50;
    }
    function y(coord) {
	    if (parseTime(coord.TimeAxis) == 0) {
		    return 50 + data.Stations[coord.SpaceAxis].SpaceAxisHeap*20;
	    }
        
	    let delta = (parseTime(coord.TimeAxis)-parseTime(data.MinTime))/1000;
	    return delta/data.TimeAxisDistance*(data.TimeAxisSize-100)+100;
    }
    function lz(i) {
        return i < 10 ? '0'+i : ''+i;
    }
    function simpleTime(t) {
        const d = new Date(t);
        return lz(d.getHours())+':'+lz(d.getMinutes());
    }

    let currentSelected = undefined;
    let currentSelectedShortestPath = [];

    function selectEdge(edgeId) {
        setSelectedForDependents(false);
        currentSelected = data.Edges[edgeId];
        console.log('cur', currentSelected, edgeId);
        currentSelectedShortestPath = setSelectedForDependents(true);
        console.log(currentSelectedShortestPath);
    }

    function setSelectedForDependents(selected) {
        if (!currentSelected) {
            return [];
        }
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
        }
        next = currentSelected;
        while(next.ReverseShortestPath.length > 0) {
            previous = next;
            next = data.Edges[next.ReverseShortestPath[0].EdgeID];
            setSelectedForEdge(next, selected);
            setSelectedForStationEdge(previous, next, selected);
            all.push(next);
        }
        return all;
    }

    function setSelectedForEdgeId(edgeId, selected) {
        const e = document.getElementById(edgeId);
        if (selected) {
            e.className.baseVal += " selected";
        } else {
            e.className.baseVal =  e.className.baseVal.replace(" selected", "");
        }    
    }

    function setSelectedForEdge(edge, selected) {
        if (edge.Discarded) return;
        setSelectedForEdgeId(edge.ID, selected);
    }
    
    
    function setSelectedForStationEdge(previous, next, selected) {
        if (previous.Discarded || next.Discarded) return;
        const edgeId = previous.ID+'_'+next.ID+'_station';
        if (document.getElementById(edgeId)) {
            setSelectedForEdgeId(edgeId, selected);
        }
    }

    function selectListener(evt) {
        const id = this.id.replace('-toucharea', '');
        console.log('selected ', id);
        selectEdge(id);
    }

    function label(e) {
        if (!e.Line) {
            return '';
        }
        let label = '';
        if (e.Line.Name) {
            label = e.Line.Name;
        } else {
            label = e.Line.ID;
        }
        if (e.Message) {
            label += ' 🛈';
        }
        if (e.Line.Type == 'Foot') {
            return '🚶 ' + label;
        }
        return label;
    }

    function type(e) {
        if (!e.Line) {
            return '';
        }
        return e.Line.Type;
    }

    function departure(e) {
        return time(e, stop => stop.Departure, stop => stop.DepartureTrack);
    }

    function arrival(e) {
        return time(e, stop => stop.Arrival, stop => stop.ArrivalTrack);
    }

    function liveDataDeparture(e) {
        return liveDataClass(e, stop => stop.Departure);
    }

    function liveDataArrival(e) {
        return liveDataClass(e, stop => stop.Arrival);
    }

    function time(e, timeResolver, trackResolver) {
        if (!e.Line) {
            return ''
        }
        let label = simpleTime(timeResolver(e.Actual)) + ' ' + delay(timeResolver(e.Current), timeResolver(e.Planned))
        if (trackResolver(e.Planned)) {
            label += "Pl." + trackResolver(e.Planned)
        }
        return label
    }

    function liveDataClass(e, timeResolver) {
        if (!e.Line) {
            return '';
        }
        let current = timeResolver(e.Current)
        if (parseTime(current) == 0) {
            return ''
        }
        if (delayMinutes(current, timeResolver(e.Planned)) > 5) {
            return "live-red"
        }
        return "live-green"
    }

    function delayMinutes(current, planned) {
        return Math.round((parseTime(current)-parseTime(planned))/1000/60);
    }

    function delay(current, planned) {
        if (parseTime(current) != 0) {

            return " (+" + delayMinutes(current, planned) + ") ";
        }
        return ''
    }

    const margin = 25;

    onMount(() => {
        fetchTimespace();
    })
</script>

<div class="loading-screen" style="display: {!data ? 'block' : 'none'};">
    <img src="res/icon/loading.gif" id="loading-indicator">
    <p>Data retrieval can take up to one minute.</p>
</div>

{#if data}
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" 
    viewBox="0 0 {data.SpaceAxisSize} {data.TimeAxisSize}">
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
{#each Object.values(data.Stations) as s}
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
{#each data.SortedEdges.map(id => data.Edges[id]) as e}
{#if !e.Discarded}
<path id="{e.ID}" d="M {x(e.From)},{y(e.From)} L{x(e.To)},{y(e.To)}"
    class="edge type-{type(e)} redundant-{e.Redundant} {e.ShortestPathFor.map(p => 'sp-'+p).join(' ')} {e.ProviderShortestPath ? 'provider-shortest-path' : ''}"
    />
<path id="{e.ID}-toucharea" d="M {x(e.From)},{y(e.From)} L{x(e.To)},{y(e.To)}"
    class="edge-toucharea" on:click={selectListener}
    />
{/if}
{/each}
{#each currentSelectedShortestPath as e}
{#if !e.Discarded}
<text id="{e.ID}-label" class="label type-{type(e)} label-{e.ID}">
    <textPath href="#{e.ID}" startOffset="50%">
        {label(e)}
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
{#if e.PreviousArrival}<text x="{x(e.To)}" y="{y(e.To)-margin}" class="previous-next-arrow" on:click={selectEdge(e.PreviousArrival)}>▲</text>{/if}
{#if e.PreviousDeparture}<text x="{x(e.From)}" y="{y(e.From)-margin}" class="previous-next-arrow" on:click={selectEdge(e.PreviousDeparture)}>▲</text>{/if}
{#if e.NextArrival}<text x="{x(e.To)}" y="{y(e.To)+margin}" class="previous-next-arrow" on:click={selectEdge(e.NextArrival)}>▼</text>{/if}
{#if e.NextDeparture}<text x="{x(e.From)}" y="{y(e.From)+margin}" class="previous-next-arrow" on:click={selectEdge(e.NextDeparture)}>▼</text>{/if}
{/if}
{/each}
<a xlink:href="?{optionsQueryString(store)}&amp;form">
    <rect x="5" y="1430" width="100" height="40" class="button" />
    <text x="55" y="1450" class="ui-link">
        Modify query
    </text>
</a>
<text x="100" y="1450" id="details">
    
</text>
<g id="legend">
    <path d="M 150,1460 L350,1460" class="edge selected redundant-false" />
    <text x="150" y="1450">
        Fastest route (transfer time >= 0 min)
    </text>
    <path d="M 400,1460 L600,1460" class="edge provider-shortest-path redundant-false" />
    <text x="400" y="1450">
        Route recommended by DB/HAFAS
    </text>
    <path d="M 650,1460 L850,1460" class="edge redundant-true" />
    <text x="650" y="1450">
        Connections deemed redundant
    </text>
    <path d="M 900,1460 L980,1460" class="edge type-nationalExpress redundant-false" />
    <text x="900" y="1450">
        ICE, IC, etc.
    </text>
    <path d="M 1030,1460 L1110,1460" class="edge type-regional redundant-false" />
    <text x="1030" y="1450">
        RE, RB, S, etc.
    </text>
    <path d="M 1160,1460 L1240,1460" class="edge type-bus redundant-false" />
    <text x="1160" y="1450">
        Bus, Tram, etc.
    </text>
    <path d="M 1290,1460 L1370,1460" class="edge type-Foot redundant-false" />
    <text x="1290" y="1450">
        On Foot
    </text>
</g>
</svg>
{/if}