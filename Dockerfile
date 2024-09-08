FROM golang:1.23-alpine

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the gRPC server
RUN go build -o server .

# Expose the gRPC port
EXPOSE 50051

# Start the server
CMD ["./server"]

