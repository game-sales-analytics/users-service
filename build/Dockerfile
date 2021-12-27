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
COPY --from=build /app/bin/ping /ping

HEALTHCHECK --interval=5s --timeout=2s --start-period=5s --retries=2 CMD [ "/ping", "userssrv:50050" ]

EXPOSE ${PORT}

USER nonroot:nonroot

ENTRYPOINT [ "/userssrv" ]
