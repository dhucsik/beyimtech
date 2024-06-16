FROM golang:1.21.5-alpine as builder

COPY . /app
WORKDIR /app

RUN apk add --no-cache gcc musl-dev libwebp-dev
RUN CGO_ENABLED=1 go build -o ./bin/service ./cmd/main.go

FROM alpine:latest as runner

RUN apk add --no-cache tzdata libwebp-dev

COPY --from=builder /app/bin/service /app/bin/service
COPY --from=builder /app/configs/app.toml /app/configs/app.toml
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/.env /app/.env

RUN mkdir /app/public

WORKDIR /app
CMD ["bin/service"]