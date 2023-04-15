function queryString(q) {
    return Object.keys(q)
    .map(
        k => Array.isArray(q[k]) 
        ? q[k].map(v => ({k: k, v: v}))
        : [{k: k, v: q[k]}]
    )
    .flat()
    .map(o => o.k+'='+(o.v ? o.v : ''))
    .join('&');
}

function optionsQueryString(query) {
    const q = {
        from: query.from.id,
        to: query.to.id,
        vias: query.vias.map(v => v.id),
        datetime: query.datetime
    };
    return queryString(q);        
}

export {
    optionsQueryString
}