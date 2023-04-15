export class Station {
    id: string | undefined;
    name: string | undefined;
}

const from = new Station();
const to = new Station();
const vias: Station[] = [];

let datetime: string | null = null;


const store = {
    from: from,
    to: to,
    vias: vias,
    datetime: datetime,
    initialized: false
};

function mapStation(s) {
    return {id: s.ID, name: s.Name};
}

function setFromApi(data) {
    store.from = mapStation(data.From);
    store.to = mapStation(data.To);
    store.vias = fillupStations(data.Vias.map(mapStation));
}

function requiredFieldsSet() {
    return store.vias.length > 0 && store.from.id && store.to.id
}

const maxVias = 10;
function fillupStations(vias) {
    const l = vias.length;
    for(let i=0; i<maxVias-l; i++) {
        vias.push(new Station());
    }
    return vias;
}

export {
    store,
    requiredFieldsSet,
    fillupStations,
    setFromApi
}