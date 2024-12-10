curl https://mirror.traines.eu/hafas-ibnr-zhv-gtfs-osm-matching/hafas-stations.ndjson | jq --slurp '. | map( { (.id): [.location.longitude, .location.latitude] } ) | add' > location-lookup.ign.json

