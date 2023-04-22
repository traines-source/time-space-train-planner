#!/bin/sh
set -e

mkdir -p static/res/fonts/
cd static/res/fonts/
wget -N https://traines.eu/res/fonts/roboto-condensed-v25-latin-700.woff
wget -N https://traines.eu/res/fonts/roboto-condensed-v25-latin-700.woff2
wget -N https://traines.eu/res/fonts/roboto-condensed-v25-latin-regular.woff
wget -N https://traines.eu/res/fonts/roboto-condensed-v25-latin-regular.woff2
wget -N https://traines.eu/res/fonts/material-icons-outlined-v109-latin-regular.woff
wget -N https://traines.eu/res/fonts/material-icons-outlined-v109-latin-regular.woff2

cd ../../
wget -N https://traines.eu/favicon.ico
wget -N https://traines.eu/apple-touch-icon.png
wget -N https://traines.eu/favicon-32x32.png
wget -N https://traines.eu/favicon-16x16.png