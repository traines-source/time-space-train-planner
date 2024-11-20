<script lang="ts">
    import { t } from '$lib/translations';
    import AutoComplete from "simple-svelte-autocomplete"
    import type { StationLabel } from '../store';

    export let selectedStation: StationLabel;
    export let placeholder: string;
    export let clearButton: boolean = true;

    function getItems(input: string) {
        return fetch(import.meta.env.VITE_STATIONS_API + "?addresses=false&poi=false&pretty=false&query="+input)
        .then(response => response.json())
        .then(data => {
            var list = data
            .filter((station: any) => !station.isMeta)
            .map((station: any) => ({
                name: station.name,
                id: station.id
            }));
            return list;
        })
        .catch((error) => {
            console.log(error);
        });
    };

</script>

<AutoComplete searchFunction={getItems}
    delay="200"
    placeholder={placeholder}
    noResultsText="{$t('c.no_results')}"
    hideArrow={true}
    localFiltering={false}
    labelFieldName="name"
    valueFieldName="id"
    cleanUserText={false}
    minCharactersToSearch="2"
    className="station"
    inputClassName="station"
    dropdownClassName="station"
    showClear={clearButton && selectedStation?.id}
    bind:selectedItem={selectedStation}
    />

