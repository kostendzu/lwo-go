FROM golang:1.22 AS builder
ENV GO111MODULE=on CGO_ENABLED=1 GOOS=linux
RUN apt-get update && apt-get install -y gcc libc-dev
WORKDIR /lwo-go
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
RUN go build -o /build/lwo-go -a -ldflags '-linkmode external -extldflags "-static"' ./cmd

FROM alpine:3

# # Устанавливаем необходимые зависимости для работы с SQLite
# RUN apk add --no-cache sqlite-libs
COPY --from=builder /build/lwo-go /bin/lwo-go
ENV FILEPATH="./tasks.db" SERVER_ADDRESS="0.0.0.0:8080" NODE_ENV=DOCKER
ENTRYPOINT ["/bin/lwo-go"]
