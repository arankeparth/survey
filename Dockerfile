# Stage 1: Build the Go application
FROM golang:1.23

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies (this layer is cached if go.mod and go.sum are not changed)
RUN go mod download

# Build the Go app
RUN go build -o main .

# Expose port 8081 to the outside world
EXPOSE 8080 6060

# Command to run the executable
CMD ["./main"]