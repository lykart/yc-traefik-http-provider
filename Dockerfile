# Build Stage
FROM golang:1.27-alpine3.24 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 go build -o yc-provider ./cmd/yc-provider/provider.go

# Run Stage
FROM alpine:3.24

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/yc-provider .
COPY ./config.yaml /app/config.yaml

# Command to run the executable
CMD ["./yc-provider"]
