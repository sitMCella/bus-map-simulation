# Hub Backend

## Introduction

Golang application for setting up the Hub. The frontend application retrieves from the Hub the bus stop and time table configurations. The bus applications send to the Hub the bus positions.

The application stores the data in a PostgreSQL database.

## Development

### Requirements

- Golang 1.25.5

### Build Application

```sh
go build
```

### Configure Application

Create a file .env with the following content:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_NAME=busmap
DB_PASSWORD=mysecretpassword
```

### Run Application

```sh
go run .
```

### Format Code

```sh
go fmt
```
