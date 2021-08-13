FROM ubuntu:latest
RUN mkdir /app
RUN mkdir /app/migrations
COPY ./bin/dp210goapp /app
COPY ./migrations /app/migrations
