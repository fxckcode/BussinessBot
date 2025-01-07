# Start from the official Go image
FROM golang:1.23-alpine

# Set the Current Working Directory inside the container
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY ./ ./

# Build the Go app
RUN go build main.go

EXPOSE 3001

# Set environment variable for Gin mode
ENV GIN_MODE=release

# Command to run the executable
CMD ["./main"]