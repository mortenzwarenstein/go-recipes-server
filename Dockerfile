# Stage 1: Build the Go binary
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (better for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o gorecipe

# Stage 2: Create a minimal image with only the binary
FROM alpine:latest

# Add a non-root user (optional but recommended)
RUN adduser -D recipe

WORKDIR /home/gorecipe

# Copy binary from builder stage
COPY --from=builder /app/gorecipe .

# Use non-root user to run app
USER recipe

EXPOSE 8080

# Command to run the binary
ENTRYPOINT ["./gorecipe"]
