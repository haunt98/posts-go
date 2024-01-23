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

## Technical side

How do services communicate with each other?

### API

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

#### References

### Message Broker

**Second** way is by Message Broker, the most well known is Kafka.

Main idea is decoupling.

Imaging service A need to call services B, C, D, E after doing some action, but
B just died. We must handle B errors gracefully if B is not that important (not
affect main flow of A). Imaging not only B, but multi B, like B1, B2, B3, ...
Bn, this is so depressed to continue.

Message Broker is a way to detach B from A.

Dumb exaplain be like: each time A do something, A produces message to Message
Broker, than A forgets about it. Then all B1, B2 can consume A's message if they
want and do something with it, A does not know and does not need to know about
it.

### Tip

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

**Pro tip**: Use proto to define models (if you can) to take advantage of
detecting breaking changes.

### References

- [Best practices for REST API design](https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/)
  - [ZaloPay API](https://docs.zalopay.vn/v2/)
  - [stripe API](https://stripe.com/docs/api)
  - [moov API](https://docs.moov.io/api/)
- [Using Apache Kafka to process 1 trillion inter-service messages](https://blog.cloudflare.com/using-apache-kafka-to-process-1-trillion-messages/)

## Coding principle

You should know about DRY, SOLID, KISS, Design Pattern. The basic is learning
which is which when you read code. Truly understand will be knowing when to use
and when to not.

All of these above are industry standard.

The way business moving is fast, so a feature is maybe implemented today, but
gets thrown out of window tomorrow (Like A/B testing, one of them is chosen, the
other says bye). So how do we adapt? The problem is to detect, which
code/function is likely stable, resisted changing and which is likely to change.

For each service, I often split to 3 layers: handler, service, repository.

- Handler layer: Handle HTTP/GRPC/Message Broker/...
- Service layer: All rules, logic goes here.
- Repository layer: Interact with cache/databases using CRUD and some cache
  strategy.

Handler layer is likely never changed. Repository layer is rarely changed.
Service layer is changed daily, this is where I put so much time on.

The previous question can be asked in many ways:

- How to move fast without breaking things?
- How to quickly experiment new code without affecting old code?
- ...

My answer is, as Message Broker introduce concept decoupling, loosely coupled
coding. Which means, 2 functions which do not share same business can be deleted
without breaking the other.

For example, we can send noti to users using SMS, Zalo, or Noti in app (3
providers). They are all independently feature which serves same purpose: alert
user about something. What happen if we add providers or remove some? Existing
providers keep working as usual, new providers should behave properly too.

So we have send noti abstraction, which can be implement by each provider, treat
like a module (think like lego) which can be plug and play right away.

And when we do not need send noti anymore, we can delete whole of it which
includes all providers and still not affecting main flow.

### References

- [Write code that is easy to delete, not easy to extend.](https://programmingisterrible.com/post/139222674273/write-code-that-is-easy-to-delete-not-easy-to)
- [Imaginary Problems Are the Root of Bad Software](https://cerebralab.com/Imaginary_Problems_Are_the_Root_of_Bad_Software)

## Known concept

TODO: Cache strategy, async operation

## Challenge

TODO: Scale problem

## Damage control

TODO: Take care incident

## Bonus

TODO
