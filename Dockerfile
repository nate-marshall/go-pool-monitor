# Build stage
FROM golang:1.22-bookworm AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o poolmonitor ./cmd/poolmonitor

# Use a minimal image as a runtime stage
FROM debian:bookworm-slim

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the build stage
COPY --from=build /app/poolmonitor .

# Copy the .env file if you have any environment variables
COPY .env .

# Expose port if necessary (for example, if you're using a web server)
# EXPOSE 8080

# Command to run the executable
CMD ["./poolmonitor"]
