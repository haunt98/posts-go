# Backend Thinking

## Backend Role

Transform business requirements to action, which usually involves:

- Service:
  - ZaloPay use microservices architecture, mostly written using Go and Java
- API:
  - HTTP (Client-Server) and GRPC (Server-Server)
- Database/Cache/Storage/Message Broker
  - MySQL/Redis/S3/Kafka
  - CRUD
- Docs
  - Mostly design notes and diagrams which show how to implement business
    requirements

After successfully do all of that, next step is:

- Testing
  - Unit tests, Integration tests
- Observation
  - Log
  - Metrics
  - Tracing
