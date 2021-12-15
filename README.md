# Time-Space Train Planner

An SVG-based tool to visualize public transport journeys retrieved from a HAFAS system, in order to see all possible connections and possibly find faster connections. See it in action in this video: https://youtu.be/rD5iATcC9Mo

![Example Diagram](res/screenshot.png?raw=true)

Often, the HAFAS routing system will not show the fastest routes, because it deems the transfer times too short or because there are just too many possibilities. This tool will display all direct connections between a given set of stations and will help you find the fastest connection and any backup connections that might be good to know about. You can click on individual connections to show the shortest route to the destination from that connection (with a minimum transfer time of zero minutes).

TSTP is currently in an early alpha stage. I.e. it is unstable, does very little error handling and is not very user-friendly yet.

## Running TSTP

The server-side part is written in Go, there is also a small client-side JavaScript part. You will have to run your own server to use this software.

TSTP relies on an adapted HAFAS API provided by [hafas-client](https://github.com/public-transport/hafas-client), more specifically, [db-rest](https://github.com/derhuerst/db-rest). You need to run your own instance of this API adapter separately from TSTP, or use a publicly available one (e.g. https://v5.db.transport.rest/).

In addition, you should run an aggressive HTTP cache in between TSTP and this API, because TSTP itself is stateless and will repeatedly issue identical requests for the same resources while gathering data. In `deployments/nginx-cache.example.conf`, you find an example of how to configure Nginx as an HTTP cache. Please note that this configuration will cache for a very long time. You will not be able to get updated live information for a particular request.

The steps to run TSTP itself are:

1. Copy `deployments/conf.example.env` to `deployments/conf.env` and fill in the host name of your HTTP cache (`API_CACHE_HOST`) and, if your setup is not exactly mirroring db-rest, the path prefix (`HAFAS_API_CACHE_PREFIX`). If you run the cache locally (e.g. in another Docker container), you probably want to adjust the `API_CACHE_SCHEME` to `http`. The DB Open API and thus the respective environment variables are currently not used anymore.
2. Copy `res/conf.example.js` to `res/conf.js` and fill in the complete URL to the /stations endpoint of your db-rest instance (preferably cached as well). Obviously, since this will be used by the browser for autocompletion, this URL needs to be accessible from the outside (i.e. no Docker hostname).
3. Start TSTP using the `docker-compose.yaml` (add a port mapping if you need to) or, if you have installed all Go dependencies, without Docker using `start.sh`. The default port is 3000.
4. You should now be able to access TSTP at your selected port under `/tstp` in your browser.
5. If something goes wrong, consider looking into TSTP's stdout logs.
