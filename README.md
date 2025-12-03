# Time-Space Train Planner

Instead of planning a fixed public transport journey from origin to destination upfront, this tool shows the next best alternatives for the current transfer situation with their destination arrival time distribution including average arrival time, most likely arrival time and the 95th percentile of arrival time. It works like a car navigation system but for public transport: You will see later how to continue further, and the itinerary down the line will be continuously adapted to the real-time situation. By clicking/tapping on a connection, you basically board it, and the next best connections to your destination from the next relevant stop of that connection are shown.

Use it online at https://tespace.traines.eu

<img src="res/screenshot0.png?raw=true" width="500">

The goal of probabilistic routing taking into account transfer reliability and alternative continuations was inspired by the original time-space diagram showing all possible connections between a set of transfer stations. This diagram can still be shown in the UI on demand, but is usually not needed anymore, thanks to the probabilities. You can see that part of TSTP in action in this video: https://youtu.be/rD5iATcC9Mo

![Example Diagram](res/screenshot.png?raw=true)

TSTP is currently in development, but you can already use it productively (I know I do).

## Running TSTP

It is recommended to use to hosted version at https://tespace.traines.eu, but you can of course host your own instance.

TSTP itself provides the UI and middleware to connect to data providers. The probabilistic routing happens in the [stochastic-journey-strategies](https://github.com/traines-source/stochastic-journey-strategies) project. 

TSTP relies on [Friendly Public Transport Format](https://github.com/public-transport/friendly-public-transport-format) APIs to obtain timetable data, provided by [motis-fptf-client](https://github.com/motis-project/motis-fptf-client) for global coverage from [Transitous](http://transitous.org/) and alternatively [hafas-client](https://github.com/public-transport/hafas-client)/[db-vendo-client](https://github.com/public-transport/db-vendo-client) for Deutsche Bahn data. You need to run your own instance of these API adapters separately from TSTP, or use a publicly available one.

In addition, you should run an aggressive HTTP cache in between TSTP and this API, because TSTP itself is stateless and will repeatedly issue identical requests for the same resources while gathering data. In `deployments/nginx-cache.conf`, you find an example of how to configure Nginx as an HTTP cache. Please note that this configuration will cache for a very long time. You will not be able to get updated live information for a particular request. If you use the provided `docker-compose.yaml` as is, an Nginx reverse proxy using this configuration will be started alongside TSTP.

The steps to run TSTP itself are:

1. Copy `deployments/conf.example.env` to `deployments/conf.env` and fill in the hostname of your cached API (`API_CACHE_HOST`) and, if your setup is not exactly mirroring db-rest, the path prefix (`HAFAS_API_CACHE_PREFIX`). If you run the cached API on another host, you probably want to adjust the `API_CACHE_SCHEME` to `https`. The DB Open API and thus the respective environment variables are currently not used anymore.
2. Copy `res/conf.example.js` to `res/conf.js` and fill in the complete URL to the /stations endpoint of your cached API. Obviously, since this will be used by the browser for autocompletion, this URL needs to be accessible from the outside (i.e. no Docker hostname).
3. Start TSTP using the `docker-compose.yaml`. If you use the configuration as is, you should access TSTP via the reverse proxy, i.e. via the host-mapped port 8080. Or, if you have installed all Go dependencies, start TSTP without Docker using `start.sh`, in which case TSTP itself will be available under port 3000.
4. Access TSTP at your selected port under `/tstp` in your browser.
5. If something goes wrong, consider looking into TSTP's stdout logs.

This setup will use the rudimentary server-rendered UI. In order to use the full-fledged UI, you would also need to run the Svelte frontend in the ui/ directory and properly connect it to the backend.