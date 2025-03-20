# Simple http server using golang

## Running Server

### Using docker:

1. clone project.
2. exec `docker-compose up --build` in project main directory.

### Building from source:

1. install golang v1.24.1.
2. clone the project.
3. exec `go build` to build application.
4. create local copy of .env file containing environment variables

```
//.env file content.
PORT=8080
JWT_SECRET=NewSecretKey
DB_NAME=storedb
DBPath={{Your_DB_Connection_String}}
```

5. create new db in postgresql `storedb`.
6. exec `./go_http_server` to run application.
