# Paging using Redis

Make use of Redis sorted set for paging.

This is my way to bypass database for paging. It's not perfect, just purely quick hacks.

Each record has `created_at_ms`, I use this as **score** for Redis sorted set, also build a string representation of
record to store in sorted set.

What field need for query will be used for string representation, for example
`request_id|user_id|created_at_ms|updated_at_ms|status`

- For create/update record, I use **ZADD** with `created_at_ms` and string representation
- For query records, I use **ZREVRANGEBYSCORE** to limit by `created_at_ms` first then parse string representation to
  limit more

The problem is once you make a string representation, it's kinda hardcoded. If you want to add new field to filter, you
will need to migrate from old string representation to new one.
