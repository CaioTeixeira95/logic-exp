# Logic Expression

Running application

```sh
$ docker-compose up db web
```

or

```sh
$ DATABASE_URL=... go run cmd/main.go
```

Running tests

```sh
$ docker-compose up db -d
$ go test -v -race -cover ./...
```
