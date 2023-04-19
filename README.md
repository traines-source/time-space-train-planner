# Time-Space Train Planner

A browser-based tool to visualize public transport journeys from A to B, in order to see all possible connections and possibly find faster connections. See it in action in this video: https://youtu.be/rD5iATcC9Mo

Use it online at https://tespace.traines.eu

![Example Diagram](res/screenshot.png?raw=true)

Often, the HAFAS routing system will not show the fastest routes, because it deems the transfer times too short or because there are just too many possibilities. This tool will display all direct connections between a given set of stations and will help you find the fastest connection and any backup connections that might be good to know about. You can click on individual connections to show the shortest route to the destination from that connection (with a minimum transfer time of zero minutes).

TSTP is currently in development, but you can already use it productively (I know I do).

## Running TSTP

The server-side part is written in Go. It is recommended to use to hosted version at https://tespace.traines.eu, but you can of course host your own instance.

TSTP relies on an adapted HAFAS API provided by [hafas-client](https://github.com/public-transport/hafas-client), more specifically, [db-rest](https://github.com/derhuerst/db-rest). You need to run your own instance of this API adapter separately from TSTP, or use a publicly available one (e.g. https://v5.db.transport.rest/).

In addition, you should run an aggressive HTTP cache in between TSTP and this API, because TSTP itself is stateless and will repeatedly issue identical requests for the same resources while gathering data. In `deployments/nginx-cache.conf`, you find an example of how to configure Nginx as an HTTP cache. Please note that this configuration will cache for a very long time. You will not be able to get updated live information for a particular request. If you use the provided `docker-compose.yaml` as is, an Nginx reverse proxy using this configuration will be started alongside TSTP.

The steps to run TSTP itself are:

1. Copy `deployments/conf.example.env` to `deployments/conf.env` and fill in the hostname of your cached API (`API_CACHE_HOST`) and, if your setup is not exactly mirroring db-rest, the path prefix (`HAFAS_API_CACHE_PREFIX`). If you run the cached API on another host, you probably want to adjust the `API_CACHE_SCHEME` to `https`. The DB Open API and thus the respective environment variables are currently not used anymore.
2. Copy `res/conf.example.js` to `res/conf.js` and fill in the complete URL to the /stations endpoint of your cached API. Obviously, since this will be used by the browser for autocompletion, this URL needs to be accessible from the outside (i.e. no Docker hostname).
3. Start TSTP using the `docker-compose.yaml`. If you use the configuration as is, you should access TSTP via the reverse proxy, i.e. via the host-mapped port 8080. Or, if you have installed all Go dependencies, start TSTP without Docker using `start.sh`, in which case TSTP itself will be available under port 3000.
4. Access TSTP at your selected port under `/tstp` in your browser.
5. If something goes wrong, consider looking into TSTP's stdout logs.

This setup will use the rudimentary server-rendered UI. In order to use the full-fledged UI, you would also need to run the Svelte frontend in the ui/ directory and properly connect it to the backend.