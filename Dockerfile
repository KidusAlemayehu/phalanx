# Build Stage
FROM golang:1.25-alpine AS builder

# Set up working directory
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the proxy executable
# CGO_ENABLED=0 builds a statically linked binary (better for scratch/alpine)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /phalanx ./cmd/phalanx/main.go

# Final Stage
FROM alpine:3.19

# Add ca-certificates in case the proxy needs to communicate externally
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /phalanx .

# The standard port we listen on
EXPOSE 8443

# Command to run the executable
ENTRYPOINT ["./phalanx"]
