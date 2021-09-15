# Build the Go Binary.
FROM golang:1.17-alpine as build_appointment-api

RUN mkdir -p /service
WORKDIR /service
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o appointment-api 

FROM alpine:latest
ARG BUILD_DATE
ARG VCS_REF
COPY --from=build_appointment-api /service /service
WORKDIR /service
CMD ["./appointment-api"]
