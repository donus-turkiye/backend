# backend
 
## wakeup
docker-compose up --build

## stop and remove containers
docker-compose down -v --remove-orphans

## swagger
## swagger address
http://localhost:8080/swagger/index.html

### install swag
go get -u github.com/swaggo/swag/cmd/swag
### add to path
export PATH=$(go env GOPATH)/bin:$PATH
### generate swagger docs
swag init -g docs/swagger.go