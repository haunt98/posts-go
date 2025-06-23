# Reload config

This serves as design draft of reload config system

```plantuml
@startuml Reload config

skinparam defaultFontName Iosevka Term SS08

participant admin
participant other_service
participant config_service
participant storage

== Admin handle ==

admin -> config_service: set/update/delete config

config_service -> storage: set/update/delete config

== Other service handle ==

other_service -> other_service: init service

activate other_service

other_service -> storage: make connection

loop

    other_service -> storage: listen on config change

    other_service -> other_service: save config to memory

end

deactivate other_service

other_service -> other_service: do business

activate other_service

other_service -> other_service: get config

other_service -> other_service: do other business

deactivate other_service

@enduml
```

Config storage can be any key value storage or database like etcd, Consul, mySQL, ...

If storage is key value storage, maybe there is API to listen on config change. Otherwise we should create a loop to get
all config from storage for some interval, for example each 5 minute.

Each `other_service` need to get config from its memory, not hit `storage`. So there is some delay between upstream
config (config in `storage`) and downstream config (config in `other_service`), but maybe we can forgive that delay
(???).

Pros:

- Config can be dynamic, service does not need to restart to apply new config.

- Each service only keep 1 connection to `storage` to listen to config change, not hit `storage` for each request.

Cons:

- Each service has 1 more dependency, aka `storage`.
- Need to handle fallback config, incase `storage` failure.
- Delay between upstream/downstream config
