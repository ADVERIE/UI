FROM golang:1.19-alpine AS builder

# Install necessary tools
RUN apk add --no-cache protobuf git curl bash protobuf-dev

# Install Go dependencies
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

# Set the working directory
WORKDIR /src/ui-service

# Copy go.mod and go.sum
COPY go.mod go.sum* ./
RUN go mod download

# Copy the proto files
COPY proto/ proto/

# Copy the script to generate proto files
COPY generate_protos.sh ./
RUN chmod +x generate_protos.sh

# Generate proto files
RUN ./generate_protos.sh

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o ui-service ./cmd/ui

# Create a lightweight production image
FROM alpine:3.17

# Install dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary and other necessary files from the builder stage
WORKDIR /app
COPY --from=builder /src/ui-service/ui-service /app/
COPY --from=builder /src/ui-service/templates /app/templates
COPY --from=builder /src/ui-service/public /app/public

# Set up health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/ || exit 1

# Set environment variables
ENV GRPC_PORT=50053
ENV HTTP_PORT=8080

# Expose ports
EXPOSE $GRPC_PORT $HTTP_PORT

# Set the entrypoint
ENTRYPOINT ["/app/ui-service"] 