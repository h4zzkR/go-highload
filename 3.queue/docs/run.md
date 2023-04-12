go get github.com/rabbitmq/amqp091-go
go get golang.org/x/crypto/bcrypt
go get -u github.com/golang-jwt/jwt/v5

### Запуск брокера
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.11-management

