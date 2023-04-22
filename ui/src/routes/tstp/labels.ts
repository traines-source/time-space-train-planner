import { t } from '$lib/translations';
import type { Edge } from './types';


function parseTime(t: string): number {
    return Math.max(0, Date.parse(t));
}

function lz(i: number): string {
    return i < 10 ? '0'+i : ''+i;
}
function simpleTime(t: string): string {
    const d = new Date(t);
    return lz(d.getHours())+':'+lz(d.getMinutes());
}

function label(e: Edge, detail: boolean): string {
    if (!e.Line) {
        return '';
    }
    let label = e.Line.Name;
    if (e.Message && !detail) {
        label += ' <tspan class="micon">info</tspan>';
    }
    if (e.Line.Type == 'Foot'  && !detail) {
        return '<tspan class="micon">directions_walk</tspan> ' + label;
    } else if (e.Line.Type == 'Foot'  && detail) {
        return '<span class="micon">directions_walk</span> ' + label;
    }
    return label;
}

function type(e: Edge): string {
    if (!e.Line) {
        return '';
    }
    return e.Line.Type;
}

function departure(e: Edge): string {
    return time(e, (stop: any) => stop.Departure, (stop: any) => stop.DepartureTrack);
}

function arrival(e: Edge): string {
    return time(e, (stop: any) => stop.Arrival, (stop: any) => stop.ArrivalTrack);
}

function liveDataDeparture(e: Edge): string {
    return liveDataClass(e, (stop: any) => stop.Departure);
}

function liveDataArrival(e: Edge): string {
    return liveDataClass(e, (stop: any) => stop.Arrival);
}

function time(e: Edge, timeResolver: (stop: any) => string, trackResolver: (stop: any) => string) {
    if (!e.Line) {
        return ''
    }
    let label = simpleTime(timeResolver(e.Actual)) + ' ' + delay(timeResolver(e.Current), timeResolver(e.Planned))
    if (trackResolver(e.Planned)) {
        label += t.get('c.platform') + trackResolver(e.Planned)
    }
    return label
}

function liveDataClass(e: Edge, timeResolver: (stop: any) => string) {
    if (!e.Line) {
        return '';
    }
    const current = timeResolver(e.Current)
    if (parseTime(current) == 0) {
        return ''
    }
    if (delayMinutes(current, timeResolver(e.Planned)) > 5) {
        return "live-red"
    }
    return "live-green"
}

function delayMinutes(current: string, planned: string) {
    return Math.round((parseTime(current)-parseTime(planned))/1000/60);
}

function delay(current: string, planned: string) {
    if (parseTime(current) != 0) {

        return " (+" + delayMinutes(current, planned) + ") ";
    }
    return ''
}

export {
    parseTime,
    simpleTime,
    label, type, departure, arrival, liveDataDeparture, liveDataArrival
}