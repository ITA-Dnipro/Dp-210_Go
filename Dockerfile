# Builder image.
FROM golang:latest as builder
RUN mkdir /app
WORKDIR /app
RUN mkdir ./bin
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/dp210goapp ./

# Production image.
FROM ubuntu:latest
RUN mkdir /app
RUN mkdir /app/migrations
COPY --from=builder /app/bin/dp210goapp /app
COPY ./migrations /app/migrations
COPY ./config.json /app
