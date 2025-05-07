FROM golang:1.24.3-bullseye AS builder

WORKDIR /app

# Copy only necessary files for dependency resolution
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o aws-identity-validator

# Use a minimal aws image for the final container
FROM amazon/aws-cli:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/aws-identity-validator /app/

# Set the entrypoint
ENTRYPOINT ["/app/aws-identity-validator"]
