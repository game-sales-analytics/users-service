# syntax=docker/dockerfile:1

FROM docker.io/library/golang:1 AS build

WORKDIR /app

COPY go.sum go.mod ./

RUN go mod download

COPY . .

RUN make build

FROM gcr.io/distroless/base-debian11

ARG PORT=5000

WORKDIR /

COPY --from=build /app/bin/userssrv /userssrv

EXPOSE ${PORT}

USER nonroot:nonroot

ENTRYPOINT [ "/userssrv" ]
