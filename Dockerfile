# Step 1: Build the Go application (builder stage)
FROM golang:tip-alpine3.22 AS builder

# Set the image metadata with labels
LABEL maintainer="tranaduc9x@gmail.com"
LABEL version="1.0.0"

# Set environment variables for Go build
ENV DOCKER_BUILDKIT=1
WORKDIR /build

# Copy go.mod and go.sum for dependency management
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Set CGO_ENABLED to 0 for static linking and set GOOS and GOARCH for Linux amd64 architecture
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Build the Go application
RUN go build -o main ./cmd/auth-service/main.go

# Step 2: Create the minimal image (distroless stage)
FROM gcr.io/distroless/static-debian12:latest

# Set working directory in the distroless container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /build/main .

# Expose the port your application will run on
EXPOSE 9420

# Start the application
CMD ["/app/main"]
