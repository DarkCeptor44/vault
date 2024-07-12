FROM golang:1.22.5-alpine AS build-env

ENV PORT=8080

WORKDIR /go/src/vault
COPY . /go/src/vault
RUN CGO_ENABLED=0 go build .

FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates && rm -rf /var/cache/apk*
WORKDIR /app
COPY ./public ./public
COPY ./views ./views
COPY --from=build-env /go/src/vault/vault /app/vault

EXPOSE ${PORT}
CMD DOCKER=yes ./vault
