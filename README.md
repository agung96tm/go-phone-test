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

### Migrate and Run Applications
```
$ make db/migration/up
$ make run/web
$ make run/api
```
nb: run `run/web` and `run/api` with separate terminals.
