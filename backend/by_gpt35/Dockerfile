# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container's working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod vendor
RUN go mod tidy

# Copy the source code from the local host to the container's working directory
COPY . .

# Build the Go app
RUN
RUN go build -o main planningpoker

# Expose the port the app runs on
EXPOSE 80

# Command to run the app
CMD ["/app/main"]
