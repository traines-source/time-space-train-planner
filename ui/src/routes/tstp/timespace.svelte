<script>
    import { onMount, tick } from "svelte";
    import { setFromApi, Station, store } from "../store"
    import { optionsQueryString } from "../url"
    import panzoom from 'panzoom'
    
    let loading = true;
    let query = store;
    let data;
    let error;
    let currentSelected = undefined;
    let currentSelectedShortestPath = [];
    const arrowMargin = 25;

    function fetchTimespace() {
        fetch(import.meta.env.VITE_TSTP_API+'timespace?'+optionsQueryString(query))
        .then(response => response.json())
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
        .catch((error) => {
            alert('Failed request. Possibly too many requests. Try again later.');
            loading = false;
            console.log(error);
        });
    }

    function refresh() {
        error = undefined;
        loading = true;
        fetchTimespace();
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

    function selectEdge(edgeId) {
        setSelectedForDependents(false);
        currentSelected = data.Edges[edgeId];
        console.log('cur', currentSelected, edgeId);
        currentSelectedShortestPath = setSelectedForDependents(true);
        console.log(currentSelectedShortestPath);
    }

    function setActive(selected) {
        const e = document.getElementById(currentSelected.ID+'-toucharea');
        if (!e) return;
        if (selected) {
            e.className.baseVal += " active";
        } else {
            e.className.baseVal =  e.className.baseVal.replace(" active", "");
        }
    }

    function setSelectedForDependents(selected) {
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

    function setSelectedForEdgeId(edgeId, selected) {
        const e = document.getElementById(edgeId);
        if (!e) return;
        if (selected) {
            e.className.baseVal += " selected";
        } else {
            e.className.baseVal =  e.className.baseVal.replace(" selected", "").replace(" active", "");
        }    
    }

    function setSelectedForEdge(edge, selected) {
        if (!edge || edge.Discarded) return;
        setSelectedForEdgeId(edge.ID, selected);
    }
    
    
    function setSelectedForStationEdge(previous, next, selected) {
        if (!previous || !next || previous.Discarded || next.Discarded) return;
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

    function label(e, detail) {
        if (!e.Line) {
            return '';
        }
        let label = '';
        if (e.Line.Name) {
            label = e.Line.Name;
        } else {
            label = e.Line.ID;
        }
        if (e.Message && !detail) {
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



    onMount(() => {
        fetchTimespace();
        panzoom(document.getElementById('timespace-canvas'), {
            maxZoom: 7,
            minZoom: 1,
            bounds: true,
            boundsPadding: 1,
            zoomDoubleClickSpeed: 3,
            onTouch: function(e) {
                if (e.target.id == 'timespace-canvas') return true; 
            }
        })
    })
</script>

<div class="loading-screen" style="display: {!data ? 'block' : 'none'};">
    <svg class="indicator {loading ? 'loading' :''}" viewBox="-48.2015 -48.2015 400 400">
        <path d="M57.866,268.881c25.982,19.891,56.887,30.403,89.369,30.402h0.002c6.545,0,13.176-0.44,19.707-1.308 c39.055-5.187,73.754-25.272,97.702-56.557c14.571-19.033,24.367-41.513,28.329-65.01c0.689-4.084-2.064-7.954-6.148-8.643 l-19.721-3.326c-1.964-0.33-3.974,0.131-5.595,1.284c-1.621,1.153-2.717,2.902-3.048,4.864 c-3.019,17.896-10.49,35.032-21.608,49.555c-18.266,23.861-44.73,39.181-74.521,43.137c-4.994,0.664-10.061,1-15.058,1 c-24.757,0-48.317-8.019-68.137-23.191c-23.86-18.266-39.18-44.73-43.136-74.519c-3.957-29.787,3.924-59.333,22.189-83.194 c21.441-28.007,54.051-44.069,89.469-44.069c24.886,0,48.484,7.996,68.245,23.122c6.55,5.014,12.43,10.615,17.626,16.754 l-36.934-6.52c-1.956-0.347-3.973,0.101-5.604,1.241c-1.631,1.141-2.739,2.882-3.085,4.841l-3.477,19.695 c-0.72,4.079,2.003,7.969,6.081,8.689l88.63,15.647c0.434,0.077,0.869,0.114,1.304,0.114c1.528,0,3.031-0.467,4.301-1.355 c1.63-1.141,2.739-2.882,3.084-4.841l15.646-88.63c0.721-4.079-2.002-7.969-6.081-8.69l-19.695-3.477 c-4.085-0.723-7.97,2.003-8.689,6.082l-6.585,37.3c-7.387-9.162-15.87-17.463-25.248-24.642 c-25.914-19.838-56.86-30.324-89.495-30.324c-46.423,0-89.171,21.063-117.284,57.787C6.454,93.385-3.878,132.123,1.309,171.178 C6.497,210.236,26.583,244.933,57.866,268.881z"></path>
    </svg>
    <p>Data retrieval can take up to one minute.</p>
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
{#each Object.values(data.Stations) as s (s.ID)}
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
    class="edge type-{type(e)} redundant-{e.Redundant} cancelled-{e.Cancelled} {e.ShortestPathFor.map(p => 'sp-'+p).join(' ')} {e.ProviderShortestPath ? 'provider-shortest-path' : ''}"
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
{#if e.PreviousArrival}<text x="{x(e.To)}" y="{y(e.To)-arrowMargin}" class="previous-next-arrow" on:click={selectEdge(e.PreviousArrival)}>▲</text>{/if}
{#if e.PreviousDeparture}<text x="{x(e.From)}" y="{y(e.From)-arrowMargin}" class="previous-next-arrow" on:click={selectEdge(e.PreviousDeparture)}>▲</text>{/if}
{#if e.NextArrival}<text x="{x(e.To)}" y="{y(e.To)+arrowMargin}" class="previous-next-arrow" on:click={selectEdge(e.NextArrival)}>▼</text>{/if}
{#if e.NextDeparture}<text x="{x(e.From)}" y="{y(e.From)+arrowMargin}" class="previous-next-arrow" on:click={selectEdge(e.NextDeparture)}>▼</text>{/if}
{/if}
{/each}
{/if}

</svg>

</div>
<div id="details"><div>
    {#if currentSelected}
    <div class="refresh"><a href="javascript:void(0)" on:click={refresh}>
        <svg class="indicator {loading ? 'loading' :''}" viewBox="-48.2015 -48.2015 400 400">
            <path d="M57.866,268.881c25.982,19.891,56.887,30.403,89.369,30.402h0.002c6.545,0,13.176-0.44,19.707-1.308 c39.055-5.187,73.754-25.272,97.702-56.557c14.571-19.033,24.367-41.513,28.329-65.01c0.689-4.084-2.064-7.954-6.148-8.643 l-19.721-3.326c-1.964-0.33-3.974,0.131-5.595,1.284c-1.621,1.153-2.717,2.902-3.048,4.864 c-3.019,17.896-10.49,35.032-21.608,49.555c-18.266,23.861-44.73,39.181-74.521,43.137c-4.994,0.664-10.061,1-15.058,1 c-24.757,0-48.317-8.019-68.137-23.191c-23.86-18.266-39.18-44.73-43.136-74.519c-3.957-29.787,3.924-59.333,22.189-83.194 c21.441-28.007,54.051-44.069,89.469-44.069c24.886,0,48.484,7.996,68.245,23.122c6.55,5.014,12.43,10.615,17.626,16.754 l-36.934-6.52c-1.956-0.347-3.973,0.101-5.604,1.241c-1.631,1.141-2.739,2.882-3.085,4.841l-3.477,19.695 c-0.72,4.079,2.003,7.969,6.081,8.689l88.63,15.647c0.434,0.077,0.869,0.114,1.304,0.114c1.528,0,3.031-0.467,4.301-1.355 c1.63-1.141,2.739-2.882,3.084-4.841l15.646-88.63c0.721-4.079-2.002-7.969-6.081-8.69l-19.695-3.477 c-4.085-0.723-7.97,2.003-8.689,6.082l-6.585,37.3c-7.387-9.162-15.87-17.463-25.248-24.642 c-25.914-19.838-56.86-30.324-89.495-30.324c-46.423,0-89.171,21.063-117.284,57.787C6.454,93.385-3.878,132.123,1.309,171.178 C6.497,210.236,26.583,244.933,57.866,268.881z"></path>
        </svg>
    </a></div>
    <div class="train">
        <h4>
            <span class="label">{label(currentSelected, true)}</span>
            <span class="destination">
                {#if currentSelected.Line && currentSelected.Line.Direction}
                <span class="arrow">➔</span> {currentSelected.Line.Direction}
                {/if}
            </span>
        </h4>
    </div>
    <div class="arrdep">
        <span class="dep">
            <span class="{liveDataDeparture(currentSelected)}">{departure(currentSelected)}</span>
            {#if currentSelected.Cancelled}<span class="cancelled">(Cancelled)</span>{/if}
            <br />{data.Stations[currentSelected.From.SpaceAxis].Name}
        </span>
        <svg viewBox="0 0 50 10" class="miniature">
            <path d="M 10,5 L40,5" class="edge type-{type(currentSelected)} redundant-false"/>
        </svg>
        <span class="arr">
            <span class="{liveDataArrival(currentSelected)}">{arrival(currentSelected)}</span>
            {#if currentSelected.Cancelled}<span class="cancelled">(Cancelled)</span>{/if}
            <br />{data.Stations[currentSelected.To.SpaceAxis].Name}
            
        </span>
    </div>
    {#if currentSelected.Message}
    <div class="message">
        🛈 {currentSelected.Message}
    </div>
    {/if}
    {/if}

    <a href="?{optionsQueryString(store)}&form" class="submit">
            Modify query
    </a>

    <div class="legend">
        <h4>Legend</h4>
        <svg viewBox="125 1430 500 150" style="width: 100%">
            <g id="legend">
                <g>
                    <path d="M 150,1460 L350,1460" class="edge selected redundant-false" />
                    <text x="150" y="1450">
                        Fastest route (transfer time >= 0 min)
                    </text>
                </g>
                <g transform="translate(250, 0)">
                    <path d="M 150,1460 L350,1460" class="edge provider-shortest-path redundant-false" />
                    <text x="150" y="1450">
                        Route recommended by DB/HAFAS
                    </text>
                </g>
                <g transform="translate(0, 40)">
                    <path d="M 150,1460 L350,1460" class="edge redundant-true" />
                    <text x="150" y="1450">
                        Connections deemed redundant
                    </text>
                </g>
                <g transform="translate(0, 80)">
                    <path d="M 150,1460 L350,1460" class="edge type-nationalExpress redundant-false" />
                    <text x="150" y="1450">
                        ICE, IC, etc.
                    </text>
                </g>
                <g transform="translate(250, 80)">
                    <path d="M 150,1460 L350,1460" class="edge type-regional redundant-false" />
                    <text x="150" y="1450">
                        RE, RB, S, etc.
                    </text>
                </g>
                <g transform="translate(0, 120)">
                    <path d="M 150,1460 L350,1460" class="edge type-bus redundant-false" />
                    <text x="150" y="1450">
                        Bus, Tram, etc.
                    </text>
                </g>
                <g transform="translate(250, 120)">
                    <path d="M 150,1460 L350,1460" class="edge type-Foot redundant-false" />
                    <text x="150" y="1450">
                        On Foot
                    </text>
                </g>
            </g>
        </svg>
    </div>
    <p id="footer">
        <a href="https://github.com/traines-source/time-space-train-planner/issues">Report an issue</a>
        <a href="{import.meta.env.VITE_TSTP_LEGAL}">Imprint/Privacy</a>
    </p>
</div></div>