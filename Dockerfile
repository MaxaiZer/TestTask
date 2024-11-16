FROM alpine:latest AS base

FROM golang:1.23 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal/
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/build/main cmd/main.go

FROM base AS final
WORKDIR /app
RUN mkdir ./configs & mkdir ./logs
COPY configs/config.env ./configs
COPY migrations ./migrations/
COPY --from=build /app/build .
CMD ["./main"]