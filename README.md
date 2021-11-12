# Time-Space Train Planner

An SVG-based tool to visualize public transport journeys retrieved from a HAFAS system, in order to see all possible connections and possibly find faster connections.

![Example Diagram](res/screenshot.png?raw=true)

Often, the HAFAS routing system will not show the fastest routes, because it deems the transfer times too short or because there are just too many possibilities. This tool will display all direct connections between a given set of stations and will help you find the fastest connection and any backup connections that might be good to know about.

## Running TSTP

The server-side part is written in Go, there is also a small client-side JavaScript part. You will have to run your own server to use this software. 
