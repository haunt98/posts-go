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

## What else?

I think about adding prefix to ID to identify which resource that ID represents.

## Thanks

- [UUID, SERIAL OR IDENTITY COLUMNS FOR POSTGRESQL AUTO-GENERATED PRIMARY KEYS?](https://www.cybertec-postgresql.com/en/uuid-serial-or-identity-columns-for-postgresql-auto-generated-primary-keys/)
- [Identity Crisis: Sequence v. UUID as Primary Key](https://brandur.org/nanoglyphs/026-ids)
- [Generating good unique ids in Go](https://blog.kowalczyk.info/article/JyRZ/generating-good-unique-ids-in-go.html)
- [How we used Go 1.18 when designing our Identifiers](https://encore.dev/blog/go-1.18-generic-identifiers)
- [ULIDs and Primary Keys](https://blog.daveallie.com/ulid-primary-keys)
- [On IDs](https://0pointer.net/blog/projects/ids.html)
