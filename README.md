# go-microservices-suite
A Go-based microservices system featuring gRPC, RabbitMQ, authentication, logging, messaging, and Kubernetes deployment.


# docker compose

```
docker-compose up
docker compose up -d --build broker-service
docker compose up -d --build frontend-service
```

# broker-service

```
go get github.com/go-chi/chi/v5
go get github.com/go-chi/chi/v5/middleware
go get github.com/go-chi/cors
```

```
docker build -f .\broker-service\broker-service.dockerfile -t broker-service:latest .\broker-service
docker run -d --name broker-service -p 8080:80 broker-service:latest
```

# frontend-service

```
docker build -f .\frontend\frontend-service.dockerfile -t frontend:latest .\frontend
docker run -d --name frontend -p 8081:80 frontend:latest
```

# authentication-service

```
go get golang.org/x/crypto/bcrypt
go get github.com/go-chi/chi/v5
go get github.com/go-chi/chi/v5/middleware
go get github.com/go-chi/cors

go get github.com/jackc/pgconn
go get github.com/jackc/pgx/v4
go get github.com/jackc/pgx/v4/stdlib
```


# logger service (to view db data inside mongo shell)
```
docker exec -it 3da31bc9fb84 mongo -u admin -p password --authenticationDatabase admin
```
```
show dbs
use logs
show collections
```
```
db.logs.find().limit(10).pretty()
```