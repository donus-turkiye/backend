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