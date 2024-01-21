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

In ZaloPay, each team has its own responsibilites/domains, aka many different
services.

Ideally each team can choose custom backend techstack if they want, but mostly
boils down to Java or Go. Some teams use Python for scripting, data processing,
...

_Example_: UM (Team User Management Core) has 10+ Java services and 30+ Go
services.

The question is for each new business requirements, what should we do:

- Create new services with new APIs?
- Add new APIs to existing services?
- Update existing APIs?
- Change configs?
- Don't do anything?

_Example_: Business requirements says: Must match/compare user EKYC data with
Bank data (name, dob, id, ...). TODO

TODO: How to split services?

## Technically side

How do services communicate with each other?

**First** is through API, this is the direct way, you send a request then you
wait for response.

**HTTP**: GET/POST/...

_Example_: TODO use curl

**GRPC**: use proto file as constract.

_Example_: TODO: show sample proto file

There are no hard rules on how to design APIs, only some best practices, like
REST API, ...

Correct answer will always be: "It depends". Depends on:

- Your audience (android, ios, web client or another internal service)
- Your purpose (allow to do what?)
- Your current techstack (technology limitation?)
- Your team (bias, prefer, ...?)
- ...

Why do we use HTTP for Client-Server and GRPC for Server-Server?

- HTTP for Client-Server is pretty standard. Easy for client to debug, ...
- Before ZaloPay switch to GRPC for Server-Server, we use HTTP. The reason for
  switch is mainly performance.

Take a look at [stripe API](https://stripe.com/docs/api) to get quick grasp
about real world API.

My own experiences:

- Whatever you design, you stick with it consistently. Don't use different name
  for same object/value in your APIs.
- Don't trust client blindly, everything can be fake, everything must be
  validated. We can not know the request is actually from our client or some
  hacker computer. (Actually we can but this is out of scope, and require lots
  of advance work)
- Don't delete/rename/change old fields because you want and you can, please
  think it through before do it. Because back compability is very hard, old apps
  should continue to function if user don't upgrade. Even if we rollout new
  version, it takes time.

## Bonus

TODO

## References

- [Best practices for REST API design](https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/)
