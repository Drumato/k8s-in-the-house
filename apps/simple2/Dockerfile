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

RUN go build -o /simple2

##
## Deploy
##
FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /simple2 /simple2

EXPOSE 1323

USER nonroot:nonroot

ENTRYPOINT ["/simple2"]

