# Start with the official Go image as the base
FROM golang:1.22-alpine AS builder

# Set necessary environment variables for Go and install Git
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    PORT=8080 \
    MONGODB_URL=mongodb+srv://MoZaki:thisIsASecret@organizationmanager.m1dpzkn.mongodb.net/?retryWrites=true&w=majority \
    SECRET_KEY=thisIsASecret

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd

# Start a new stage from scratch
FROM alpine:latest  

# Set necessary environment variables
ENV PORT=$PORT \
    MONGODB_URL=$MONGODB_URL \
    SECRET_KEY=$SECRET_KEY

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /app/main

# Expose port 8080 to the outside world
EXPOSE $PORT

# Command to run the executable
CMD ["/app/main"]