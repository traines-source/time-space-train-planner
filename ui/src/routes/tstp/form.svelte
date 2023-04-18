<script>
    import { Station, store, fillupStations, setFromApi } from "../store"
    import { optionsQueryString } from "../url"
    import StationInput from "./stationInput.svelte"
    import { goto } from '$app/navigation';
    import { onMount } from "svelte";

    let query = store;
    let loading = false;
   
    function fetchVias() {
        fetch(import.meta.env.VITE_TSTP_API+'vias?'+optionsQueryString(query))
        .then(response => response.json())
        .then(data => {
            console.log(data);
            setFromApi(data);
            query = query;
            loading = false;
        })
        .catch((error) => {
            alert('Failed request. Possibly too many requests. Try again later.');
            loading = false;
            console.log(error);
        });
    }

    function submit() {
        if (query.from.id && query.to.id) {
            loading = true;
            goto('?'+optionsQueryString(query));
            if (query.vias.length == 0) {
                fetchVias();
            }
        }
    }

    onMount(() => {
        if (!query.from.name || !query.to.name) {
            if (query.from.id && query.to.id) {
                loading = true;
                fetchVias();
            }
        }
    })

</script>

<div id="header"><div>
    <a href="/tstp"><h1>TeSpace<sup>BETA</sup></h1></a>
</div></div>
<div id="container">
    <p>
        This tool allows you to plan your public transport journeys using an interactive time-space diagram, based on live timetable data (currently Germany only). This is an early Beta version and as such might stop working or return wrong data at any time.    
    </p>
        
    <form autocomplete="off" id="query">

        <div>
            <StationInput placeholder="From" bind:selectedStation={query.from} />
        </div>
        <div>
            <StationInput placeholder="To" bind:selectedStation={query.to} />
        </div>
        {#if query.vias.length > 0}
            <p>
                You need to enter relevant interchange stations (up to 10). Only direct connections between these stations will be found.
            </p>
        {/if}

        {#each query.vias as via, i}
            <div>
                <StationInput placeholder="Via" bind:selectedStation={query.vias[i]} />
            </div>
        {/each}

        <div>
            <input type="datetime-local" id="datetime" name="datetime" bind:value="{query.datetime}">
            <p id="default-now">Default is Now</p>
        </div>
        <div>
            <input type="checkbox" id="regionly" bind:checked="{query.regionly}"><label for="regionly"> Regional transport only</label>
        </div>
        
        <div id="submit-container"><input type="button" value="Submit" class="submit" on:click={submit}><!--
        --><img src="res/icon/loading.gif" id="loading-indicator" style="display: {loading ? 'block' : 'none'};"></div>
        
        Data retrieval can take up to one minute. 
    </form>
</div>