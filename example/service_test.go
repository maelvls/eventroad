package service

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/maelvls/eventroad/eventpb"
	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/proto"
)

func TestBankAccount_Apply(t *testing.T) {
	testCases := map[string]struct {
		entity     BankAccount
		event      proto.Message
		wantEntity BankAccount
		wantPanic  bool
	}{
		"create": {
			entity:     BankAccount{},
			event:      &eventpb.Created{Name: "foo"},
			wantEntity: BankAccount{Name: "foo"},
		},
		"create uses an empty entity before applying": {
			entity:     BankAccount{Name: "non-fresh-entity"},
			event:      &eventpb.Created{Name: ""},
			wantEntity: BankAccount{Name: ""},
		},
		"edited": {
			entity:     BankAccount{Name: "foo"},
			event:      &eventpb.Edited{Name: "bar"},
			wantEntity: BankAccount{Name: "bar"},
		},
		"edited doesn't override with zero values": {
			entity:     BankAccount{Name: "foo"},
			event:      &eventpb.Edited{Name: ""},
			wantEntity: BankAccount{Name: "foo"},
		},
		"unknown event": {
			entity:    BankAccount{},
			event:     &eventpb.UnhandledEvent{},
			wantPanic: true,
		},
	}
	for testName, tt := range testCases {
		t.Run(testName, func(t *testing.T) {
			got := tt.entity

			if tt.wantPanic {
				require.Panics(t, func() {
					got.Apply(tt.event)
				})
				return
			}

			got.Apply(tt.event)
			assert.Equal(t, tt.wantEntity, got)
		})
	}
}
