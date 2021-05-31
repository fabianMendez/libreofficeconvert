FROM golang:1.16 as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o main .

FROM ubuntu:20.04

COPY --from=builder /build/main /app/

ENV PORT=8080
EXPOSE 8080

WORKDIR /app
ENTRYPOINT ["./main"]
