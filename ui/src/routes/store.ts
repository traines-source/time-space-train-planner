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

function requiredFieldsSet(): boolean {
    return store.vias.length > 0 && !!store.from.id && !!store.to.id
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
    fillupStations,
    setFromApi
}