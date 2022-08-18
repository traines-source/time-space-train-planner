var stationMap = {};
$(".station-autocomplete").autocomplete({
    source: function( request, response ) {
        fetch(STATIONS_API + "?addresses=false&poi=false&pretty=false&query="+request.term)
        .then(response => response.json())
        .then(data => {
            console.log(data);
            var list = data
            .filter(station => !station.isMeta)
            .map(station => {
                console.log(station.name, station["name"]);
                stationMap[station.name] = station.id;
                return station.name;
            });
            console.log(list)
            response(list);
        })
        .catch((error) => {
            alert('Failed autocomplete request. Possibly too many requests. Try again later.')
        });
    }
});

function prepareSubmit() {
    if (mapEvaNumbers()) {
        document.getElementById('loading-indicator').style.display = 'block';
        return true;
    }
    return false;
}

function mapEvaNumbers() {
    const inputs = document.getElementsByClassName('station-autocomplete');
    for (var i=0;i<inputs.length;i++) {
        if (!mapEvaNumber(inputs[i])) {
            return false;
        }
    }
    return true;
}

function mapEvaNumber(inputField) {
    const evaNumber = stationMap[inputField.value];
    const id = inputField.id.replace('-name', '');
    if (evaNumber == undefined) {
        if ((id == 'from' || id == 'to') && !inputField.value) {
            inputField.value = '';
            return false
        }
        return true
    }
    document.getElementById(id).value = evaNumber;
    return true;
}