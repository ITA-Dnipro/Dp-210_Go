# Builds binary
FROM golang as builder
LABEL stage=compilation
RUN mkdir /build
WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -tags netgo -a -o doctor-api
RUN ls

# Runs binary
FROM alpine:latest
COPY --from=builder /build/sql /sql
COPY --from=builder /build/doctor-api /
RUN ls
CMD [./doctor-api]
