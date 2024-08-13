# syntax=docker/dockerfile:1
FROM golang:1.22-alpine3.20 AS build-env

ENV PORT=8080
ENV DOCKER=yes

WORKDIR /go/src/vault
COPY . /go/src/vault
RUN CGO_ENABLED=0 go build -o ./vault .

FROM alpine:3.20
RUN apk --no-cache add ca-certificates && \
        mkdir -p /app && \
        chown -R nobody: /app

WORKDIR /app
COPY --chown=nobody:nogroup ./public ./public
COPY --chown=nobody:nogroup ./views ./views
COPY --chown=nobody:nogroup --from=build-env /go/src/vault/vault /app/vault

EXPOSE ${PORT}
ENTRYPOINT ["./vault"]
