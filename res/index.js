if (document.getElementById('datetime').value == "") {
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    now.setSeconds(0);
    now.setMilliseconds(0);
    document.getElementById('datetime').value = now.toISOString().slice(0, -8);
}

var stationMap = {};
$(".station-autocomplete").autocomplete({
    source: function( request, response ) {
        fetch(STATIONS_API + "?addresses=false&poi=false&pretty=false&query="+request.term)
        .then(response => response.json())
        .then(data => {
            console.log(data);
            var list = data.map(station => {
                console.log(station.name, station["name"]);
                stationMap[station.name] = station.id;
                return station.name;
            });
            console.log(list)
            response(list);
        });
    }
});

function prepareSubmit() {
    const inputs = document.getElementsByClassName('station-autocomplete');
    for (var i=0;i<inputs.length;i++) {
        if (!mapEvaNumber(inputs[i])) {
            return false;
        }
    }
    return true;
}

function mapEvaNumber(inputField) {
    if (inputField.disabled) {
        return true;
    }
    const evaNumber = stationMap[inputField.value];
    const id = inputField.id.replace('-name', '');
    if (evaNumber == undefined) {
        if (id == 'from' || id == 'to') {
            inputField.value = '';
            return false
        }
        return true
    }
    document.getElementById(id).value = evaNumber;
    return true;
}