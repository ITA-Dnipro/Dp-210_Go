# Builder image.
FROM golang:latest as builder
RUN mkdir /app
WORKDIR /app
RUN mkdir ./bin
COPY ./ ./
RUN go build -o ./bin/dp210goapp ./

# Production image.
FROM ubuntu:latest
RUN mkdir /app
RUN mkdir /app/migrations

RUN apt-get update -y > /dev/null
RUN apt-get install -y ca-certificates > /dev/null
RUN update-ca-certificates --fresh > /dev/null

COPY ./token.json /app
COPY ./migrations /app/migrations
COPY ./config.json /app
COPY --from=builder /app/bin/dp210goapp /app
