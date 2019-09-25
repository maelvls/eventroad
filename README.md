# eventroad

Event-sourcing library PoC with NATS

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
