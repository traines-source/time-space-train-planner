const knownErrors = [400, 500, 502, 503, 504];

function handleHttpErrors(response: Response) {
    if (!response.ok) {
        const errorMsgId = 'error_http_'+ (knownErrors.indexOf(response.status) != -1 ? response.status : 'unknown');
        throw new Error(errorMsgId);
    }
    return response.json();
}

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

function optionsQueryString(query: any, datetime?: string): any {
    const q = {
        from: query.from.id,
        to: query.to.id,
        vias: query.vias.map((v: any) => v?.id),
        datetime: query.datetime || datetime,
        regionly: query.regionly,
        system: query.system
    };
    return queryString(q);        
}

export {
    optionsQueryString,
    handleHttpErrors
}