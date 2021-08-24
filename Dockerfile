FROM golang:latest as builder
RUN mkdir /app
WORKDIR /app
RUN mkdir ./bin
COPY ./ ./
RUN go build -o ./dp210goapp ./
EXPOSE 8000
