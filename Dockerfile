# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory to /go/src/go-stock-cli
WORKDIR /app

RUN ["mkdir", "-p", "data/csv"]

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Download project dependencies
RUN go mod download

# Copy the local package files to the container's workspace
COPY stock/ stock/
COPY main.go .

# Download project dependencies
RUN ["go", "mod", "download"]
RUN ["go", "mod", "verify"]
# Build the Go app
RUN ["go", "build", "-o", "go-stock-cli"]

VOLUME ["/app/data"]

# Run the your-cli-tool command by default when the container starts
CMD ["./go-stock-cli"]
