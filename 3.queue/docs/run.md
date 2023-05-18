### Зависимости
`go get github.com/rabbitmq/amqp091-go && go get golang.org/x/crypto/bcrypt && go get -u github.com/golang-jwt/jwt/v5`

### Запуск брокера
`docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.11-management`

### Настройка почтового сервера
`export MAIL=<your-mail> && export PASSWD=<your-mail-password>`
Пример, как настроить gmail почту: https://stackoverflow.com/questions/10147455/how-to-send-an-email-with-gmail-as-provider-using-python/51664129#51664129

### Запуск сервера
`go build && ./queue`
Это сервер, диспатчит http запросы, отправляет сообщения о событиях в очередь

### Запуск consumer-email-deliver
`go run ./sender/sender.go`
Вспомогательный почтовый модуль, разгребает очередь и отправляет письма по нужным адресам



