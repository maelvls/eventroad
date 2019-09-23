package service

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	streaming "github.com/maelvls/eventroad"
	"github.com/maelvls/eventroad/example/eventpb"
)

// BankAccount is an example of an entity.
type BankAccount struct {
	Name string `json:"name"`
}

// Apply takes a BankAccount and mutates it in-place using the given event.
// If the given event doesn't have a type known by Apply, Apply will panic.
//
// Possible errors: Subject isn't know by Apply (likely) or proto.Unmarshal
// returns an error (unlikely). Apply will never error because of
// preconditions unmet regarding the entity state or the event content.
func (cur *BankAccount) Apply(s streaming.Subject, event []byte) error {
	switch s.Action {

	case "Created":
		ev := &eventpb.Created{}
		if err := proto.Unmarshal(event, ev); err != nil {
			return err
		}
		cur.Name = ev.GetName()
		return nil

	case "Edited":
		ev := &eventpb.Edited{}
		if err := proto.Unmarshal(event, ev); err != nil {
			return err
		}
		if ev.GetName() != "" {
			cur.Name = ev.GetName()
		}
		return nil

	default:
		return fmt.Errorf("Apply doesn't know how to deal with event type '%s'. Raw event: %v", event, string(event))
	}
}
