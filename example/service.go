package service

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/maelvls/eventroad/eventpb"
)

// BankAccount is an example of an entity.
type BankAccount struct {
	Name string `json:"name"`
}

// Apply takes a BankAccount and mutates it in-place using the given event.
// If the given event doesn't have a type known by Apply, Apply will panic.
func (cur *BankAccount) Apply(event proto.Message) {
	if ev, ok := event.(*eventpb.Created); ok {
		cur.Name = ev.GetName()
		return
	}
	if ev, ok := event.(*eventpb.Edited); ok {
		if ev.GetName() != "" {
			cur.Name = ev.GetName()
		}
		return
	}

	panic(fmt.Errorf("Apply doesn't know how to deal with this event: %#v", event))
}
