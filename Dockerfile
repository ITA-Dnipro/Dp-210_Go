FROM ubuntu:latest
RUN mkdir /app
RUN mkdir /app/migrations

RUN apt-get update -y > /dev/null
RUN apt-get install -y ca-certificates > /dev/null
RUN update-ca-certificates --fresh > /dev/null

COPY ./config.json /app
COPY ./migrations /app/migrations
COPY ./token.json /app
COPY ./bin/dp210goapp /app

