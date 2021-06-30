FROM golang:1.16-buster as build
ENV TZ=Asia/Tokyo

WORKDIR /src/

COPY cmd cmd
COPY config config
COPY database database
COPY internal internal
COPY plugins plugins
COPY router router
COPY shared shared
COPY transport transport
COPY utils utils
COPY main.go go.mod go.sum Makefile ./

RUN make bin/ncs-proxy/static

FROM alpine:3

RUN addgroup -g 70 -S proxy \
    && adduser -u 70 -S -D -G proxy -H -h /var/lib/proxy -s /bin/sh proxy \
    && mkdir -p /var/lib/proxy \
    && chown -R proxy:proxy /var/lib/proxy

COPY --from=build --chown=proxy:proxy /src/bin/ncs-proxy .

RUN chmod +x ncs-proxy \
    && mv ncs-proxy /usr/local/bin/ \
    && apk --no-cache --update add redis supervisor

COPY docker/supervisord.conf /etc/supervisord.conf

STOPSIGNAL SIGINT

USER proxy

CMD [ "supervisord" ]