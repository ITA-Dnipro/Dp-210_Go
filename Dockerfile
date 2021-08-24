# Builder image.
FROM golang:latest as builder
RUN mkdir /app
WORKDIR /app
RUN mkdir ./bin
COPY ./ ./
RUN go build -o ./bin/dp210goapp ./
EXPOSE 8000

# Production image.
FROM ubuntu:latest
RUN mkdir /app
RUN mkdir /app/migrations
RUN apt update
RUN apt install -y curl
RUN apt install net-tools
COPY ./migrations /app/migrations
COPY ./config.json /app
COPY --from=builder /app/bin/dp210goapp /app
EXPOSE 8000
