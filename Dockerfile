# Build stage using Golang.
FROM golang:1.22 AS builder

WORKDIR /app

# Download module dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o urlshortener-api .

# Final stage using a minimal image.
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/urlshortener-api .

EXPOSE 8080

CMD ["./urlshortener-api"]
