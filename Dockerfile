# Use an official Golang runtime as a parent image
FROM golang:1.21.1

# Set the working directory in the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Download all the dependencies
RUN go mod download

# Install the package
RUN go build -o /app/api ./cmd/api
# This container exposes port 4000 to the outside world
EXPOSE 4000

# Run the binary program produced by `go install`
CMD ["/app/api"]
