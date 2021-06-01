FROM caddy

RUN wget https://github.com/spothq/cryptocurrency-icons/archive/refs/tags/v0.17.2.zip -O icons.zip

RUN apk add --no-cache unzip
RUN unzip icons.zip -d /opt/icons

COPY cryptoicons.Caddyfile Caddyfile

CMD caddy run
