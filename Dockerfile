FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build --ldflags "-s -w -extldflags -static" -o main .

FROM alpine:latest

RUN apk add --no-cache wget

WORKDIR /www

COPY --from=builder /build/main /www/
COPY --from=builder /build/.env /www/.env
COPY --from=builder /build/public/ /www/public/
COPY --from=builder /build/resources/ /www/resources/

# Hits GET /healthcheck (unauthenticated, unversioned — unlike /ping) on the
# port baked into .env's APP_PORT.
HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=5 \
  CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:3000/healthcheck || exit 1

ENTRYPOINT ["/www/main"]
