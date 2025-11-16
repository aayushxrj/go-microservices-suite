# go-microservices-suite
A Go-based microservices system featuring gRPC, RabbitMQ, authentication, logging, messaging, and Kubernetes deployment.


# docker compose

```
docker-compose up
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