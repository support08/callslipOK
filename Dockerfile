# Use an official Go runtime as a parent image
FROM golang:1.22

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the container's WORKDIR
COPY . .

# Download any needed modules
RUN go mod download

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
