FROM ubuntu:latest
RUN mkdir /app
RUN mkdir /app/migrations
COPY ./config.json /app
COPY ./migrations /app/migrations
COPY ./token.json /app
COPY ./bin/dp210goapp /app

