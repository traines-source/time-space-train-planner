function queryString(q: any): string {
    return Object.keys(q)
    .map(
        k => Array.isArray(q[k]) 
        ? q[k].map((v: any) => ({k: k, v: v}))
        : [{k: k, v: q[k]}]
    )
    .flat()
    .map(o => o.k+'='+(o.v ? o.v : ''))
    .join('&');
}

function optionsQueryString(query: any): any {
    const q = {
        from: query.from.id,
        to: query.to.id,
        vias: query.vias.map((v: any) => v.id),
        datetime: query.datetime,
        regionly: query.regionly
    };
    return queryString(q);        
}

export {
    optionsQueryString
}