Phone Go

## How to Run

### Run DB
```
$ docker-compose up -d --build
```

### Install golang-migrate
```
$ curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
$ mv migrate.linux-amd64 $GOPATH/bin/migrate
$ migrate -version
4.14.1
```

### Run migrations
```
$ export DB_DSN="user=phone_user dbname=phone_db password=phone_password sslmode=disable"
$ migrate -path=./migrations -database=$DB_DSN up
```

### Run Projects
```
$ go run ./cmd/api
$ go run ./cmd/web
```