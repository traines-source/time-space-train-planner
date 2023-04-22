FROM node AS builder

COPY . /tmp/working/
RUN cd /tmp/working/ && npm install
RUN cd /tmp/working/ && sh ./download.sh && npm run build

FROM nginx

COPY --from=builder /tmp/working/build/ /usr/share/nginx/html/

COPY nginx.conf /etc/nginx/conf.d/default.conf
