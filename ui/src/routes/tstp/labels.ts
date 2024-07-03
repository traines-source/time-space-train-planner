import { t } from '$lib/translations';
import type { Edge } from './types';


function parseTime(t: string): number {
    return Math.max(0, Date.parse(t));
}

function lz(i: number): string {
    return i < 10 ? '0'+i : ''+i;
}
function simpleTime(t: string | number): string {
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
    if (e.Line.Type == 'Foot' && !detail) {
        return '<tspan class="micon">directions_walk</tspan> ' + label;
    } else if (e.Line.Type == 'Foot'  && detail) {
        return '<span class="micon">directions_walk</span>&nbsp;' + label;
    }
    return label;
}

function type(e: Edge): string {
    if (!e.Line) {
        return '';
    }
    return e.Line.Type;
}

function departure(e: Edge, span='span'): string {
    return timeString(e, (stop: any) => stop.Departure, (stop: any) => stop.DepartureTrack, span);
}

function arrival(e: Edge, span='span'): string {
    return timeString(e, (stop: any) => stop.Arrival, (stop: any) => stop.ArrivalTrack, span);
}

function makeSpan(span: string, clazz: string, innerHtml: string) {
    return '<'+span+' class="' + clazz + '">' + innerHtml + '</'+span+'>';
}

function timeString(e: Edge, timeResolver: (stop: any) => string, trackResolver: (stop: any) => string, span: string) {
    if (!e.Line) {
        return '';
    }
    const timeLabel = simpleTime(timeResolver(e.Actual)) + delay(timeResolver(e.Current), timeResolver(e.Planned), e);
    let label = makeSpan(span, liveDataClass(e, timeResolver), timeLabel);
    if (trackResolver(e.Actual)) {
        const platform_label = e.Line.Type == "bus" || e.Line.Type == "tram" ? t.get('c.bus_platform') : t.get('c.platform');
        label += ' <br />' + makeSpan(span, trackChangedClass(e, trackResolver), platform_label + trackResolver(e.Actual).replace(' ', '&nbsp;'));
    }
    return label;
}

function liveDataClass(e: Edge, timeResolver: (stop: any) => string) {
    if (!e.Line || e.Line.Type == 'Foot') {
        return '';
    }
    const current = timeResolver(e.Current);
    if (parseTime(current) == 0) {
        return ''
    }
    if (delayMinutes(current, timeResolver(e.Planned)) > 5) {
        return 'live-red';
    }
    return 'live-green';
}

function trackChangedClass(e: Edge, trackResolver: (stop: any) => string) {
    if (!e.Line) {
        return '';
    }
    if (trackResolver(e.Actual) != trackResolver(e.Planned)) {
        return 'live-red';
    }
    return '';
}

function delayMinutes(current: string, planned: string) {
    return Math.round((parseTime(current)-parseTime(planned))/1000/60);
}

function delay(current: string, planned: string, e: Edge) {
    if (parseTime(current) != 0 && e.Line.Type != 'Foot') {
        let m = delayMinutes(current, planned);
        return "&nbsp;(" + (m >= 0 ? '+' : '') + m + ")";
    }
    return ''
}

export {
    parseTime,
    simpleTime,
    label, type, departure, arrival
}