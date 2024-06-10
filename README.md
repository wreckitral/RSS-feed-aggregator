# RSS Feed Aggregator
The RSS Feed Aggregator API is a high-performance service built using Go, 
using sqlc for type-safe SQL queries and goose for database migrations. 
This API allows you to manage RSS feeds and articles.

## How to Run
This API can be run on your local development system using two methods.

### Directly
if you have golang and postgresql installed
- set your env 
```
PORT=<port>
DBCONN=postgres://<username>:<password>@<hostname>:5432/<dbname>?sslmode=disable
```
- install [goose](https://github.com/pressly/goose)
- `cd sql/schema`
- `goose postgres "postgres://<username>:<password>@<hostname>:5432/<dbname>?sslmode=disable" up`
- `make run`

### Using Docker
to run the project,
- `docker-compose up -d`
to turn down the project
- `docker-compose down`

## Endpoints 

