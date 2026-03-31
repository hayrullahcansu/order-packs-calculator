## Build
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o /order-packs-calculator ./src/cmd/api

## Runtime
FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates sqlite-libs

COPY --from=builder /order-packs-calculator .
COPY --from=builder /app/src/cmd/api/static ./static

ENTRYPOINT ["./order-packs-calculator"]
