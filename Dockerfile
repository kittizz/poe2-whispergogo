# Application runtime as a Builder
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .



RUN go mod download
COPY . .

ENV GOCACHE=/root/.cache/go-build

RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 go build -o __main__ cmd/server/main.go

# Application operator
FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/__main__ .

EXPOSE 80

CMD ["./__main__"]