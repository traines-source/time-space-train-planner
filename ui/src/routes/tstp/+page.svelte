<script>
    import { browser } from '$app/environment'; 
    import { requiredFieldsSet, fillupStations, store } from "../store"
    import { page } from '$app/stores'
    import Form from './form.svelte'
    import Timespace from './timespace.svelte'

    if (browser && !store.initialized) {
        store.initialized = true;
        store.from.id = $page.url.searchParams.get('from') || undefined;
        store.to.id = $page.url.searchParams.get('to') || undefined;
        const vias = $page.url.searchParams.getAll('vias').map(s => ({id: s, name: undefined}));
        if (vias.length > 0) store.vias = fillupStations(vias);
        store.datetime = $page.url.searchParams.get('datetime');
        store.regionly = $page.url.searchParams.get('regionly') == 'true';
    }

    let showForm = false;
    let showTimespace = false;
    $: {
        showForm = !browser || $page.url.searchParams.has('form') || !$page.url.searchParams.has('vias') || !requiredFieldsSet();
        showTimespace = !showForm;
    }
</script>

{#if showForm}
    <Form />
{:else if showTimespace}
    <Timespace />
{/if}