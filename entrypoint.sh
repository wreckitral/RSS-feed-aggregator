#!/bin/sh
set -e

# Wait for PostgreSQL to be available
until pg_isready -h db -p 5432 -U postgres; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done

# Run Goose migrations
goose -dir ./sql/schema postgres "postgres://postgres:rss-feed@db:5432/blogator?sslmode=disable" up

# Start the application
exec ./rss-feed
