export class StationLabel {
    id: string | undefined;
    name: string | undefined;
}

class Store {
    from = new StationLabel();
    to = new StationLabel();
    vias: StationLabel[] = [];
    datetime: string | null = null;
    regionly = false;
    system: string | null = null;
    initialized = false;
}

const store = new Store();
const defaultDatetime = getDatetime();

function getDatetime() {
    const date = new Date();
    return date.getFullYear() +
      '-' + pad(date.getMonth() + 1) +
      '-' + pad(date.getDate()) +
      'T' + pad(date.getHours()) +
      ':00';
}

function pad(num: number) {
    return (num < 10 ? '0' : '') + num;
}

function mapStation(s: any): StationLabel {
    return {id: s.ID, name: s.Name};
}

function setFromApi(data: any): void {
    store.from = mapStation(data.From);
    store.to = mapStation(data.To);
    store.vias = fillupStations(data.Vias.map(mapStation));
}

function viasSet(): boolean {
    return store.vias.filter(v => v?.id).length > 0;
}

function requiredFieldsSet(): boolean {
    return viasSet() && !!store.from.id && !!store.to.id
}

const maxVias = 10;
function fillupStations(vias: StationLabel[]): StationLabel[] {
    const l = vias.length;
    for(let i=0; i<maxVias-l; i++) {
        vias.push(new StationLabel());
    }
    return vias;
}

export {
    store,
    defaultDatetime,
    requiredFieldsSet,
    viasSet,
    fillupStations,
    setFromApi
}