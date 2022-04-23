# Start from the latest golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./gateway/gateway/go.mod ./gateway/gateway/go.sum ./

# Copy the local dependency
COPY /util ../util

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY ./gateway/gateway/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .



######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy certificates
COPY --from=builder /app/certificates ./certificates

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./main"]