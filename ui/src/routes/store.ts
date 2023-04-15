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
    fillupStations
}