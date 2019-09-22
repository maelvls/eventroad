package streaming

import "github.com/golang/protobuf/proto"

// ApplyableEntity wraps Apply. Apply applies the given event to the
// entity. Apply mutates the receiver.
//
// When given an event that isn't known by the Apply function, Apply will
// panic since it is considered a programmer mistake.
type ApplyableEntity interface {
	Apply(event interface{})
}

// ApplyFunc returns a new entity that results of the application of the
// `event` on top of the `entity`.
type ApplyFunc func(entity interface{}, event proto.Message) interface{}

// Subject represents a stream of events that you can subscribe from and is
// defined by a name such as "bankaccount" (it is called 'subject' in NATS
// and 'topic' in Rabbitmq). Each event focuses on one entity ID, and an
// entity is the 'rehydration' (applying one by one) of all events that
// have the same entity ID.
//
// Examples: BankAccount.OTE3Yzk3Y2YtMDg.Created where OTE3Yzk3Y2YtMDg is
// the id of a specific bank account.
type Subject string

// Server allows you to interact with a streaming server such as NATS or
// Kafka.
type Server interface {

	// RehydrateEntity rehydrates an entity by replaying events that match
	// the given `entityID` and store the result in the pointer `entity`.
	//
	// Subject example: BankAccount.OTE3Yzk3Y2YtMDg.* will rehydrate the
	// entity OTE3Yzk3Y2YtMDg.
	RehydrateEntity(s Subject, entity ApplyableEntity) error

	// PublishEvent publishes the given `event` to the server. The event is
	// applied to the given `entityID`.
	//
	// Possible errors are server connectivity or proto.Marshal error.
	//
	// Subject example: BankAccount.OTE3Yzk3Y2YtMDg.Created will publish
	// `event` as a Created event.
	PublishEvent(s Subject, event proto.Message) error
}
