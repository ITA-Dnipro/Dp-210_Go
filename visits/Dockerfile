# Build the Go Binary.
FROM golang:1.16-alpine as build_visits-api

RUN mkdir -p /service
WORKDIR /service
COPY . .

WORKDIR /service/visits/cmd
RUN go build -o visits-api 

FROM alpine:latest
ARG BUILD_DATE
ARG VCS_REF
COPY --from=build_visits-api /service/cmd /service
WORKDIR /service
CMD ["./visits-api"]
