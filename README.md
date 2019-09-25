# eventroad

Event-sourcing library PoC with NATS

1. how do we make sure that two published events at the same time won't conflict with each
   other? For example, two producers could process two commands on the same subject at the
   same time:
   ```
   producer1: e1 := rehydrate "id-1 #981"
   producer2: e2 := rehydrate "id-1 #981"
   producer1: publish event (entity is now "id-1 #982")
   producer2: publish event: "id-1" is now inconsistant
   ```
   But this issue is very rare; it requires the two producers to emit about the same
   entity at the exact same time since the delay between `rehydrate` and `publish` is
   very short (I mean, I hope it is).
   
2. how do we deal with old events (projection-only concern). That adds a lot of complexity
   to any consumer; if we have 3 consumers of a stream, the migration of old events (provided
   they are versionned e.g. with a `version:9` field) then the logic of migration will have to
   be written and maintained in three different places
   
   
