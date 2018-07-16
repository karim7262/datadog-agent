FROM alpine:3.8

LABEL maintainer "Datadog <package@datadoghq.com>"

RUN apk --no-cache add ca-certificates wget \
 && wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub \
 && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.27-r0/glibc-2.27-r0.apk \
 && apk add --no-cache glibc-2.27-r0.apk \
 && rm glibc-2.27-r0.apk \
 && apk del wget

COPY dogstatsd /dogstatsd

EXPOSE 8125/udp

CMD ["/dogstatsd", "start"]
