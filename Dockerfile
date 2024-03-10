# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go app
RUN ["go", "build", "-o", "go-stock-price-prediction"]

# Run the your-cli-tool command by default when the container starts
CMD ["./go-stock-price-prediction"]
