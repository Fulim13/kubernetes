# Use the official Golang image as the base image
FROM golang:latest

# Define variables
ARG port=5678
ARG project=blog

# Define environment variables
ENV project=$project
ENV port=$port

# Create a directory inside the container to store the application code
RUN mkdir /app

# Set the working directory to /app
WORKDIR /app

# Copy all files from the current directory to /app
COPY . .

# Build the Go application
RUN go build -o myapp .

# Expose the application port
EXPOSE $port

# Run the application
CMD ["./myapp"]
