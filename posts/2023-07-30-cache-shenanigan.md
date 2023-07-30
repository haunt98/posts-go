# Cache shenanigan

My notes/mistakes/... when using cache (mainly Redis) from time to time

My default strategy is:

- Write to database first then to cache second
- Read from cache first, if not found then read from database second, then re-write to cache

```mermaid
sequenceDiagram
    participant other
    participant service
    participant cache
    participant db

    other ->> service: get data
    activate service
    service ->> cache: get data
    alt exist in cache
        service -->> other: return data
    else not exist in cache
        service ->> db: get data
        alt exist data in db
            service -->> other: return data
        else not exist data in db
            service -->> other: return error not found
        end
    end
    deactivate service

    other ->> service: set data
    activate service
    service ->> db: set data
    service ->> cache: set data
    deactivate service
```

It's good for general cases, for example with CRUD action.

The bad things happen when cache and database are not consistent.
For example what happen if writing database OK then writing cache failed?
Now database has new value, but cache has old value
Then when we read again, we read cache first with old value, and that is disaster.

## Thanks

- [Cache Consistency with Database](https://danielw.cn/cache-consistency-with-database)
