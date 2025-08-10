<script>
    import { t, locales, locale } from '$lib/translations';     
</script>

<p id="footer">
    {#each $locales as value}
        <a href="javascript:void(0)" on:click="{() => $locale = value}">{value.toUpperCase()}</a> |
        <span></span>
    {/each}
    <a href="https://github.com/traines-source/time-space-train-planner/issues">{$t('c.report')}</a> |
    <a href="{import.meta.env.VITE_TSTP_LEGAL}">{$t('c.imprint')}</a> | <a target="_blank" href="{import.meta.env.VITE_DATA_ATTRIBUTION}">{$t('c.source')}</a> 
</p>