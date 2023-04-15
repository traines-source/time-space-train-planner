<script>
    import { requiredFieldsSet, fillupStations, store } from "../store"
    import { page } from '$app/stores'
    import Form from './form.svelte'
    import Timespace from './timespace.svelte'

    if (!store.initialized) {
        store.initialized = true;
        store.from.id = $page.url.searchParams.get('from') || undefined;
        store.to.id = $page.url.searchParams.get('to') || undefined;
        store.vias = fillupStations($page.url.searchParams.getAll('vias').map(s => ({id: s, name: undefined})));
    }

    let showForm;
    $: showForm = $page.url.searchParams.has('form') || !requiredFieldsSet() ;   
</script>

{#if showForm}
    <Form />
{:else}
    <Timespace />
{/if}