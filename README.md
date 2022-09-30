# Golang - Tax Calculator

Receives information from a program that consumes messages from RabbitMQ.

### Folders
- cmd: contains the main modules
- internal: contains business rules
- pkg: contains the RabbitMQ lib

### Goapp on Docker
```
docker-compose up -d
docker-compose exec goapp bash
```

### Producer
```
go run cmd/producer/main.go
```

### Consumer
```
go run cmd/consumer/main.go
```

### Grafana
localhost:3000

### RabbitMQ
localhost:15672
