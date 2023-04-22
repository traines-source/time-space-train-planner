<script lang="ts">
    import { t } from '$lib/translations';
    import { store } from "../store"
    import { optionsQueryString } from "../url"
    import Footer from '../footer.svelte'; 
    import {label, type, departure, arrival, liveDataDeparture, liveDataArrival} from './labels';

    export let currentSelected: any;
    export let loading: boolean;
    export let doRefresh: () => void;
    export let stationResolver: (id: string) => any;
</script>

<div id="details"><div>
    {#if currentSelected}
    <div class="refresh"><a href="javascript:void(0)" on:click={doRefresh}>
        <svg class="indicator {loading ? 'loading' :''}" viewBox="-48.2015 -48.2015 400 400">
            <path d="M57.866,268.881c25.982,19.891,56.887,30.403,89.369,30.402h0.002c6.545,0,13.176-0.44,19.707-1.308 c39.055-5.187,73.754-25.272,97.702-56.557c14.571-19.033,24.367-41.513,28.329-65.01c0.689-4.084-2.064-7.954-6.148-8.643 l-19.721-3.326c-1.964-0.33-3.974,0.131-5.595,1.284c-1.621,1.153-2.717,2.902-3.048,4.864 c-3.019,17.896-10.49,35.032-21.608,49.555c-18.266,23.861-44.73,39.181-74.521,43.137c-4.994,0.664-10.061,1-15.058,1 c-24.757,0-48.317-8.019-68.137-23.191c-23.86-18.266-39.18-44.73-43.136-74.519c-3.957-29.787,3.924-59.333,22.189-83.194 c21.441-28.007,54.051-44.069,89.469-44.069c24.886,0,48.484,7.996,68.245,23.122c6.55,5.014,12.43,10.615,17.626,16.754 l-36.934-6.52c-1.956-0.347-3.973,0.101-5.604,1.241c-1.631,1.141-2.739,2.882-3.085,4.841l-3.477,19.695 c-0.72,4.079,2.003,7.969,6.081,8.689l88.63,15.647c0.434,0.077,0.869,0.114,1.304,0.114c1.528,0,3.031-0.467,4.301-1.355 c1.63-1.141,2.739-2.882,3.084-4.841l15.646-88.63c0.721-4.079-2.002-7.969-6.081-8.69l-19.695-3.477 c-4.085-0.723-7.97,2.003-8.689,6.082l-6.585,37.3c-7.387-9.162-15.87-17.463-25.248-24.642 c-25.914-19.838-56.86-30.324-89.495-30.324c-46.423,0-89.171,21.063-117.284,57.787C6.454,93.385-3.878,132.123,1.309,171.178 C6.497,210.236,26.583,244.933,57.866,268.881z"></path>
        </svg>
    </a></div>
    <div class="train">
        <h4>
            <span class="label">{label(currentSelected, true)}</span>
            <span class="destination">
                {#if currentSelected.Line && currentSelected.Line.Direction}
                <span class="arrow">➔</span> {currentSelected.Line.Direction}
                {/if}
            </span>
        </h4>
    </div>
    <div class="arrdep">
        <span class="dep">
            <span class="{liveDataDeparture(currentSelected)}">{departure(currentSelected)}</span>
            {#if currentSelected.Cancelled}<span class="cancelled">({$t('c.cancelled')})</span>{/if}
            <br />{stationResolver(currentSelected.From.SpaceAxis).Name}
        </span>
        <svg viewBox="0 0 50 10" class="miniature">
            <path d="M 10,5 L40,5" class="edge type-{type(currentSelected)} redundant-false"/>
        </svg>
        <span class="arr">
            <span class="{liveDataArrival(currentSelected)}">{arrival(currentSelected)}</span>
            {#if currentSelected.Cancelled}<span class="cancelled">({$t('c.cancelled')})</span>{/if}
            <br />{stationResolver(currentSelected.To.SpaceAxis).Name}
            
        </span>
    </div>
    {#if currentSelected.Message}
    <div class="message">
        🛈 {currentSelected.Message}
    </div>
    {/if}
    {/if}

    <a href="?{optionsQueryString(store)}&form" class="submit">
        {$t('c.modify_query')}
    </a>

    <div class="legend">
        <h4>{$t('c.legend')}</h4>
        <svg viewBox="125 1430 500 150" style="width: 100%">
            <g id="legend">
                <g>
                    <path d="M 150,1460 L350,1460" class="edge selected redundant-false" />
                    <text x="150" y="1450">
                        {$t('c.fastest_route')}
                    </text>
                </g>
                <g transform="translate(250, 0)">
                    <path d="M 150,1460 L350,1460" class="edge provider-shortest-path redundant-false" />
                    <text x="150" y="1450">
                        {$t('c.provider_recommended_route')}
                    </text>
                </g>
                <g transform="translate(0, 40)">
                    <path d="M 150,1460 L350,1460" class="edge redundant-true" />
                    <text x="150" y="1450">
                        {$t('c.redundant_connection')}
                    </text>
                </g>
                <g transform="translate(250, 40)">
                    <path d="M 150,1460 L350,1460" class="edge redundant-true cancelled-true" />
                    <text x="150" y="1450">
                        {$t('c.cancelled_connection')}
                    </text>
                </g>
                <g transform="translate(0, 80)">
                    <path d="M 150,1460 L350,1460" class="edge type-nationalExpress redundant-false" />
                    <text x="150" y="1450">
                        {$t('c.long_distance')}
                    </text>
                </g>
                <g transform="translate(250, 80)">
                    <path d="M 150,1460 L350,1460" class="edge type-regional redundant-false" />
                    <text x="150" y="1450">
                        {$t('c.regional')}
                    </text>
                </g>
                <g transform="translate(0, 120)">
                    <path d="M 150,1460 L350,1460" class="edge type-bus redundant-false" />
                    <text x="150" y="1450">
                        {$t('c.local')}
                    </text>
                </g>
                <g transform="translate(250, 120)">
                    <path d="M 150,1460 L350,1460" class="edge type-Foot redundant-false" />
                    <text x="150" y="1450">
                        {$t('c.pedestrian')}
                    </text>
                </g>
            </g>
        </svg>
    </div>
    <Footer />
</div></div>