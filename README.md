# go-gin-mongodb-test

A very basic RESTful API using [Gin](https://github.com/gin-gonic/gin) and MongoDB.

There's little to no request validation. Attempting to get an object using an invalid id will result in an error.

There is a sibling repo that uses [gorilla/mux](https://github.com/gorilla/mux) instead of gin: https://github.com/ismaelpadilla/go-mux-mongodb-test

## Defaults
* The API runs on port 8080. Try navigating to http://localhost:8080/test
* Mongo Express runs on port 8081.
* The default mongo db url is: `mongodb://root:root@mongo:27017/`

## How to run

Build the go image from the `go` folder:

```sh
docker build --tag gin-test .
```

From the main folder, run everything with

```sh
docker-compose up -d
```

## Using the API

### Get all ojbects

```sh
curl 'http://localhost:8080/stuff'
```

### Get object by ID

```sh
curl 'http://localhost:8080/stuff/621be48a568a2b625835fef2'
```

### Insert object

```sh
curl -X POST 'http://localhost:8080/stuff' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Title": "A nice title",
    "Body": "This sure is a body"
}'
```

### Delete object

```sh
curl -X DELETE 'http://localhost:8080/stuff/621be48a568a2b625835fef41'
```
