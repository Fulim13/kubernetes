# Use the official Golang image as the base image
FROM golang:latest

# Define variables
ARG port=5679
ARG project=lottery

# Define environment variables
ENV project=$project
ENV port=$port

# Create a directory inside the container to store our app code
RUN mkdir /app

# Switch the working directory to /app
WORKDIR /app

# Copy all files from the current directory into the /app directory
COPY . .

# Build the Go application
RUN go build -o myapp .

# Expose the port
EXPOSE $port

# Run the application
CMD ["./myapp"]
