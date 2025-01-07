# Stage 1: Build the Go application
FROM golang:1.23.0 AS builder

# Set the working directory inside the container
WORKDIR /hrapplication

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hrapplication .

# Stage 2: Create the minimal runtime image
FROM alpine:latest

# Install required CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory in the runtime container
WORKDIR /root/

# Copy the built binary from the builder image
COPY --from=builder /hrapplication .

# Expose the necessary port (default is 443)
EXPOSE 443

# Command to run the application
CMD ["./hrapplication"]
