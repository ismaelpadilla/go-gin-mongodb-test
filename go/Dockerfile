# syntax=docker/dockerfile:1

# BUILD
FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /gin-test

# DEPLOY
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /gin-test /gin-test

USER nonroot:nonroot

ENTRYPOINT [ "/gin-test" ]
