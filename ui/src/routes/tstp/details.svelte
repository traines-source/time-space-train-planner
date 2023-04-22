<script lang="ts">
    import { t } from '$lib/translations';
    import { store } from "../store"
    import { optionsQueryString } from "../query"
    import Footer from '../footer.svelte'; 
    import {label, type, departure, arrival, liveDataDeparture, liveDataArrival} from './labels';
    import type { Edge, Station } from './types';

    export let currentSelected: Edge;
    export let loading: boolean;
    export let doRefresh: () => void;
    export let stationResolver: (id: string) => Station;
    export let error: string | undefined;
</script>

<div id="details"><div>
    {#if error}
    <p class="error">{$t('c.error')}: {$t('c.'+error)}</p>
    {/if}
    {#if currentSelected}
    <div class="refresh"><a href="javascript:void(0)" on:click={doRefresh}>
        <span class="indicator {loading ? 'loading' :''}"><span class="micon">autorenew</span></span>
    </a></div>
    <div class="train">
        <h4>
            <span class="label">{@html label(currentSelected, true)}</span>
            <span class="destination">
                {#if currentSelected.Line && currentSelected.Line.Direction}
                <span class="micon">east</span> {currentSelected.Line.Direction}
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
        <span class="micon">info</span> {currentSelected.Message}
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