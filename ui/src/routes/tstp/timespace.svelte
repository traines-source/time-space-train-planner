<script>
    import { onMount } from "svelte";
    import { Station, store } from "../store"
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
    class="edge type-{e.Type} redundant-{e.Redundant} {e.ShortestPathFor.map(p => 'sp-'+p).join(' ')} {e.ProviderShortestPath ? 'provider-shortest-path' : ''}"
    />
<path id="{e.ID}-toucharea" d="M {x(e.From)},{y(e.From)} L{x(e.To)},{y(e.To)}"
    class="edge-toucharea"
    />
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