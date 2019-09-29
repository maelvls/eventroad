# eventroad

Event-sourcing library PoC with NATS

## Terminology

- a **subject** = a string of the form `BankAccount.Ly8gRmV0Y2g.Created`.

- a **projection service** = a service that subscribes to a subject such as
  `BankAccount.*.*` in order to provide a 'view' of the data.
- a **command service** = a service that (1) subscribes temporarily to a
  subject `BankAccount.Ly8gRmV0Y2g.*` for rehydrating the entity at a given
  time in order to check that the business rules are met before applying
  the command and (2) publish an event if that command meets the business
  rules.

## Two producers emitting two conflicting events

**tl;dr:** one producer with one execution queue should be fine (they
didn't need it for LMAX so we should be fine). If we really want that (but
we don't), then we can use the solution (see below).

How do we make sure that two published events at the same time won't
conflict with each other? For example, two producers could process two
commands on the same subject at the same time:

```plain
producer1: e1 := rehydrate "id-1 #981"
producer2: e2 := rehydrate "id-1 #981"
producer1: publish event (entity is now "id-1 #982")
producer2: publish event: "id-1" is now inconsistant
```

But this issue is very rare; it requires the two producers to emit about
the same entity at the exact same time since the delay between `rehydrate`
and `publish` is very short (I mean, I hope it is).

Martin Fowler's LMAX:

> The disruptors I've described are used in a style with one producer and
> multiple consumers, but this isn't a limitation of the design of the
> disruptor. The disruptor can work with multiple producers too, in this
> case it still doesn't need locks, although it does need to use CAS
> instructions (Compare And Swap, e.g. “lock cmpxchg” on x86, see
> [lmax-paper][]).

```plain
producer1: e1 := rehydrate "id-1"         seq(id-1)=981
producer2: e2 := rehydrate "id-1"         seq(id-1)=981
producer1: publish event for "id-1"       seq(id-1)=982
producer2: seq(id-1) is not 981 -> try again
producer2: e2 := rehydrate "id-1"         seq(id-1)=982
producer1: publish event for "id-1"       seq(id-1)=983
```

[lmax-paper]: https://lmax-exchange.github.io/disruptor/files/Disruptor-1.0.pdf

## Migration of old events

How do we deal with old events (projection-only concern). That adds a lot
of complexity to any consumer; if we have 3 consumers of a stream, the
migration of old events (provided they are versionned e.g. with a
`version:9` field) then the logic of migration will have to be written
and maintained in three different places

## Existing ES projects

Our project mixes ES, CQRS and PubSub.

- https://github.com/mantzas/incata: ES with relational DBs, does not takle
  the CQRS side but has an interesting `Appender` interface. Event is

  ```go
    // Event this is the main event that will get written
    type Event struct {
        ID        int64             # the ID given by the DB
        SourceID  uuid.UUID         # the ID we give ourselves
        Created   time.Time
        Payload   interface{}
        EventType string
        Version   int
    }
  ```

- https://github.com/pavelnikolov/eventsourcing-go: not very interesting
  since it uses a in-mem `map[string]interface{}` DB. Only the
  `event-sourcing` branch has the notions of events and broker.

- https://github.com/botchniaque/eventsourcing-cqrs-go uses in-mem DB; it
  goes full Command/Aggregate but nothing really interesting

- https://github.com/savaki/eventsource has the same idea as [1]:

  ```go
  type Event interface {
      AggregateID() string
      EventVersion() int
      EventAt() time.Time
  }
  ```

  and a base impl that all events should embed:

  ```Go
  type Model struct {
      ID string
      Version int
      At time.Time
  }
  ```

  The event store interface looks like this (the store is implemented as an
  in-mem map):

  ```go
    type Store interface {
        // Save saves events to the store
        Save(aggregateID string, records ...Record) error
        // Fetch retrieves the History of events with the specified aggregate id
        Fetch(aggregateID string, version int) (History, error)
    }
  ```
