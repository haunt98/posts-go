# Dynamic config

```mermaid
sequenceDiagram
    participant admin_tool
    participant service
    participant storage

    admin_tool ->> storage: create/update/delete config

    service ->> service: init

    loop Get config periodically
        service ->> storage: get config
        service ->> service: save config to internal memory
    end

    service->>service: get config in internal memory
```

I choose Redis as config storage.

- Must have fallback value if can not get config from storage
- There is a small delay between updating config and service get the new config, because service update config
  periodically, not real time.
