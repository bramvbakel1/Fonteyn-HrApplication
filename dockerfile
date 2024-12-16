# Use official Golang image as base
FROM golang:1.23.0 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o hrapplication .

# Start a new stage from a smaller image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the builder stage
COPY --from=builder /app/hrapplication .

# Expose port 80
EXPOSE 8080

# Command to run the executable
CMD ["./hrapplication"]
