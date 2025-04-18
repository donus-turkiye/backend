# backend
 
## wakeup
docker-compose up --build

## stop and remove containers
docker-compose down -v --remove-orphans

## swagger
### install swag
go get -u github.com/swaggo/swag/cmd/swag
### add to path
export PATH=$(go env GOPATH)/bin:$PATH
### generate swagger docs
swag init -g docs/swagger.go