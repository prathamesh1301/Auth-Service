FROM golang:1.25-alpine AS build

WORKDIR /app

# Copy go mod and sum files if you have go.sum
COPY go.mod go.sum* ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o /main ./cmd

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /root/

# Copy the pre-built binary
COPY --from=build /main .
# Copy migrations so they are available (if the app runs them later)
COPY --from=build /app/migrations ./migrations

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
