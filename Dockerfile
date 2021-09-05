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
COPY ./migrations /app/migrations
COPY ./config.json /app
COPY --from=builder /app/bin/dp210goapp /app
