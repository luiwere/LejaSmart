FROM golang:1.25-bullseye AS builder
WORKDIR /src

# cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source and build with CGO enabled for sqlite3
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /app/server

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates libsqlite3-0 && rm -rf /var/lib/apt/lists/*
WORKDIR /app

# copy binary and static assets
COPY --from=builder /app/server .
COPY static ./static
COPY templates ./templates

EXPOSE 8080
ENV PORT=8080
CMD ["/app/server"]
