<script lang="ts">
    import { t } from '$lib/translations';    
    import { store, setFromApi, defaultDatetime, viasSet } from "../store"
    import { handleHttpErrors, optionsQueryString } from "../query"
    import StationInput from "./stationInput.svelte"
    import { goto } from '$app/navigation';
    import { onMount } from "svelte";

    let query = store;
    let loading = false;
    let error: string | undefined;
   
    function fetchVias(): void {
        fetch(import.meta.env.VITE_TSTP_API+'vias?'+optionsQueryString(query, defaultDatetime))
        .then(handleHttpErrors)
        .then(data => {
            setFromApi(data);
            query = query;
            loading = false;
        })
        .catch((err) => {
            loading = false;
            console.log('Error:', err);
            error = err.message && err.message.startsWith('error_http_') ? err.message : 'error_unknown';
        });
    }

    function submit(): void {
        if (query.from.id && query.to.id) {
            loading = true;
            goto('?'+optionsQueryString(query));
            if (!viasSet()) {
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
        {$t('c.introduction')}
    </p>

    {#if error}
    <p class="error">{$t('c.error')}: {$t('c.'+error)}</p>
    {/if}
        
    <form autocomplete="off" id="query">

        <div>
            <StationInput placeholder="{$t('c.from')}" bind:selectedStation={query.from} />
        </div>
        <div>
            <StationInput placeholder="{$t('c.to')}" bind:selectedStation={query.to} />
        </div>
        {#if query.vias.length > 0}
            <p>
                {$t('c.interchanges_explanation')}
            </p>
        {/if}

        {#each query.vias as via, i}
            <div>
                <StationInput placeholder="{$t('c.via')}" bind:selectedStation={query.vias[i]} />
            </div>
        {/each}

        <div>
            <input type="datetime-local" id="datetime" name="datetime" bind:value="{query.datetime}">
            <p id="default-now">{$t('c.default_now')}</p>
        </div>
        <div>
            <input type="checkbox" id="regionly" bind:checked="{query.regionly}"><label for="regionly"> {$t('c.regionly')}</label>
        </div>
        
        <div id="submit-container"><input type="button" value="{$t('c.submit')}" class="submit" on:click={submit}><!--
        --><span class="indicator {loading ? 'loading' :''}" style="visibility: {loading ? 'visible' : 'hidden'};"><span class="micon">autorenew</span></span></div>
        
        {$t('c.data_retrieval_waiting')}
    </form>
</div>