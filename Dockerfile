FROM golang:1.22.4-alpine

WORKDIR /app

# Install necessary packages including Goose and netcat
RUN apk add --no-cache git postgresql-client \
    && go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application code
COPY . ./

# Build the Go application
RUN go build -o rss-feed .

# Make the application executable
RUN chmod +x ./rss-feed

# Create log directory
RUN mkdir -p /var/log/rss-feed

# Copy the entrypoint script
COPY entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

# Expose the application's port
EXPOSE 7777

# Use the entrypoint script
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

