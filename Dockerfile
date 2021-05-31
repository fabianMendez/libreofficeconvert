FROM golang:1.16 as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:3.13

RUN \
    apk add --no-cache --virtual=build-dependencies && \
    apk add --no-cache libreoffice openjdk8-jre && \
    apk del --purge build-dependencies

COPY --from=builder /build/main /bin/main

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["main"]
