FROM golang:1.16-alpine as build_users-api

RUN mkdir -p /service
WORKDIR /service
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

WORKDIR /service/cmd
RUN go build -o users-api 

FROM alpine:latest
ARG BUILD_DATE
ARG VCS_REF
COPY --from=build_users-api /service/cmd /service
COPY --from=build_users-api /service/migrations /service/migrations
WORKDIR /service
CMD ["./users-api"]
