# Step 1: Build the Go application (builder stage)
FROM golang@sha256:ac67716dd016429be8d4c2c53a248d7bcdf06d34127d3dc451bda6aa5a87bc06 AS builder

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
RUN go build -o main ./cmd/app/main.go

# Step 2: Create the minimal image (distroless stage)
FROM gcr.io/distroless/static-debian12@sha256:c0f429e16b13e583da7e5a6ec20dd656d325d88e6819cafe0adb0828976529dc

# Set working directory in the distroless container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /build/main .

# Expose the port your application will run on
EXPOSE 9420

# Start the application
CMD ["/app/main"]
