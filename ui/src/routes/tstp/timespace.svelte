<script>
    import { onMount, tick } from "svelte";
    import { setFromApi, Station, store } from "../store"
    import { optionsQueryString } from "../url"
    import panzoom from 'panzoom'
    
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
            tick().then(() => selectEdge(data.DefaultShortestPathID));
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

    function setActive(selected) {
        const e = document.getElementById(currentSelected.ID+'-toucharea');
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
            e.className.baseVal =  e.className.baseVal.replace(" selected", "").replace(" active", "");
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

    const margin = 25;

    onMount(() => {
        fetchTimespace();
        panzoom(document.getElementById('timespace-canvas'), {
            maxZoom: 5,
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
    <img src="res/icon/loading.gif" id="loading-indicator">
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
{/if}

</svg>

</div>
<div id="details"><div>
    {#if currentSelected}
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
        <span class="dep"><span class="{liveDataDeparture(currentSelected)}">{departure(currentSelected)}</span><br />{data.Stations[currentSelected.From.SpaceAxis].Name}</span>
        <svg viewBox="0 0 50 10" class="miniature">
            <path d="M 10,5 L40,5" class="edge type-{type(currentSelected)} redundant-false"/>
        </svg>
        <span class="arr"><span class="{liveDataArrival(currentSelected)}">{arrival(currentSelected)}</span><br />{data.Stations[currentSelected.To.SpaceAxis].Name}</span>
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