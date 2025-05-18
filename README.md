# backend
 
## wakeup
docker-compose up --build

## stop and remove containers
docker-compose down -v --remove-orphans

## swagger
## swagger address
http://localhost:8080/swagger/index.html

### install swag
``` bash
go get -u github.com/swaggo/swag/cmd/swag
```
### add to path
``` bash
export PATH=$(go env GOPATH)/bin:$PATH
```
### generate swagger docs
``` bash
swag init -g docs/swagger.go
```

## migrate
### install migrate
``` bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
### add to path
``` bash
export PATH=$PATH:$(go env GOPATH)/bin
```
### create migration
``` bash
migrate -path ./migrations \
  -database "postgres://<user>:<password>@<host>:5432/<dbname>?sslmode=require" \
  up
```