import { t } from '$lib/translations';


function parseTime(t) {
    return Math.max(0, Date.parse(t));
}

function lz(i) {
    return i < 10 ? '0'+i : ''+i;
}
function simpleTime(t) {
    const d = new Date(t);
    return lz(d.getHours())+':'+lz(d.getMinutes());
}

function label(e, detail) {
    if (!e.Line) {
        return '';
    }
    let label = '';
    if (e.Line.Name) {
        label = e.Line.Name;
    } else {
        label = e.Line.ID;
    }
    if (e.Message && !detail) {
        label += ' 🛈';
    }
    if (e.Line.Type == 'Foot') {
        return '🚶 ' + label;
    }
    return label;
}

function type(e) {
    if (!e.Line) {
        return '';
    }
    return e.Line.Type;
}

function departure(e) {
    return time(e, stop => stop.Departure, stop => stop.DepartureTrack);
}

function arrival(e) {
    return time(e, stop => stop.Arrival, stop => stop.ArrivalTrack);
}

function liveDataDeparture(e) {
    return liveDataClass(e, stop => stop.Departure);
}

function liveDataArrival(e) {
    return liveDataClass(e, stop => stop.Arrival);
}

function time(e, timeResolver, trackResolver) {
    if (!e.Line) {
        return ''
    }
    let label = simpleTime(timeResolver(e.Actual)) + ' ' + delay(timeResolver(e.Current), timeResolver(e.Planned))
    if (trackResolver(e.Planned)) {
        label += t.get('c.platform') + trackResolver(e.Planned)
    }
    return label
}

function liveDataClass(e, timeResolver) {
    if (!e.Line) {
        return '';
    }
    let current = timeResolver(e.Current)
    if (parseTime(current) == 0) {
        return ''
    }
    if (delayMinutes(current, timeResolver(e.Planned)) > 5) {
        return "live-red"
    }
    return "live-green"
}

function delayMinutes(current, planned) {
    return Math.round((parseTime(current)-parseTime(planned))/1000/60);
}

function delay(current, planned) {
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