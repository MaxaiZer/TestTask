FROM alpine:latest AS base

FROM golang:1.23 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src
COPY docs/ ./docs
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/build/main src/cmd/main.go

FROM base AS final
WORKDIR /app
COPY ./config.yaml .
COPY --from=build /app/build .
COPY --from=build /app/docs .
CMD ["./main"]