FROM debian:10-slim

WORKDIR /opt/app

COPY build/server .
COPY .env .
COPY settings.toml .

STOPSIGNAL SIGINT

ENTRYPOINT ["/opt/app/server"]