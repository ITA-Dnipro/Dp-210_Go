# Build the Go Binary.
FROM golang:1.17-alpine as build_user-api

RUN mkdir -p /service
WORKDIR /service
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o user-api 

FROM alpine:latest
ARG BUILD_DATE
ARG VCS_REF
COPY --from=build_user-api /service /service
WORKDIR /service
CMD ["./user-api"]
