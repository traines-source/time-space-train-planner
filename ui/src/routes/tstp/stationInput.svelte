<script>
    import AutoComplete from "simple-svelte-autocomplete"

    export let selectedStation;
    export let placeholder;

    function getItems(input) {
        return fetch(import.meta.env.VITE_STATIONS_API + "?addresses=false&poi=false&pretty=false&query="+input)
        .then(response => response.json())
        .then(data => {
            var list = data
            .filter(station => !station.isMeta)
            .map(station => ({
                name: station.name,
                id: station.id
            }));
            return list;
        })
        .catch((error) => {
            alert('Failed autocomplete request. Possibly too many requests. Try again later.')
        });
    };

</script>

<AutoComplete searchFunction={getItems}
    delay="200"
    placeholder={placeholder}
    noResultsText="No results"
    hideArrow={true}
    localFiltering={false}
    labelFieldName="name"
    valueFieldName="id"
    cleanUserText={false}
    minCharactersToSearch="2"
    className="station"
    inputClassName="station"
    dropdownClassName="station"
    bind:selectedItem={selectedStation}
    />

<!--
   
-->

