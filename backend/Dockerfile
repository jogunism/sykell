FROM golang:1.24-alpine

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the application source code (will be volume mounted in dev)
COPY . .

# Expose the port the application listens on
EXPOSE 8080

# Command to run the executable (will be overridden by docker-compose for dev)
CMD ["go", "run", "main.go"]
