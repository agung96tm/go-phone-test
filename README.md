Go Phone Test
=========================================

### How to Run

#### 1. Run DB
```
$ docker-compose up -d --build
```
nb: depend on your prefer, you can run postgres without docker-compose.

#### 2. Copy Env
```
$ cp .envrc.example .envrc
```
nb: Please fill existing variables to run application properly.

#### 3. Install golang-migrate
```
$ curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
$ mv migrate.linux-amd64 $GOPATH/bin/migrate
$ migrate -version
4.14.1
```

#### 4. Migrate and Run Applications
```
$ make db/migration/up
$ make run/web
$ make run/api
```
nb: run `run/web` and `run/api` with separate terminals.

#### 5. Enjoy ☕
* Web: http://localhost:3000/login
* Backend: http://localhost:8000

------------

### Contributors
* Agung Yuliyanto: [Github](https://github.com/agung96tm), [LinkedIn](https://www.linkedin.com/in/agung96tm/)
