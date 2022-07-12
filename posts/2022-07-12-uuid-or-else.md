# UUID or else

There are many use cases where we need to use a unique ID.
In my experience, I only encouter 2 cases:

- ID to trace request from client to server, from service to service (microservice architecture or nanoservice I don't know).
- Primary key for database.

In my Go universe, there are some libs to help us with this:

- [google/uuid](https://github.com/google/uuid)
- [rs/xid](https://github.com/rs/xid)
- [segmentio/ksuid](https://github.com/segmentio/ksuid)
- [oklog/ulid](https://github.com/oklog/ulid)

## First use case is trace ID, or context aware ID

The ID is used only for trace and log.
If same ID is generated twice (because maybe the possibilty is too small but not 0), honestly I don't care.
When I use that ID to search log , if it pops more than things I care for, it is still no harm to me.

My choice for this use case is **rs/xid**.
Because it is small (not span too much on log line) and copy friendly.

## Second use case is primary key, also hard choice

Why I don't use auto increment key for primary key?
The answer is simple, I don't want to write database specific SQL.
SQLite has some different syntax from MySQL, and PostgreSQL and so on.
Every logic I can move to application layer from database layer, I will.

In the past and present, I use **google/uuid**, specificially I use UUID v4.
In the future I will look to use **segmentio/ksuid** and **oklog/ulid** (trial and error of course).
Both are sortable, but **google/uuid** is not.
The reason I'm afraid because the database is sensitive subject, and I need more testing and battle test proof to trust those libs.
