# Build and run

`
cd service
docker-compose up --remove-orphans
`

# Benchmark

`
ab -k -c 350 -n 20000 127.0.0.1/date
`

Условия: серверная часть работает в контейнере, запросы отправляются с помощью ab на той же машине, но вне контейнера.