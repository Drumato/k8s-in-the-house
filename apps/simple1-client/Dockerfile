# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.20.3-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go .

RUN go build -o /client

##
## Deploy
##
FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /client /client

USER nonroot:nonroot

ENTRYPOINT ["/client"]
