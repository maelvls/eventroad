package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	streaming "github.com/maelvls/eventroad"
	"github.com/maelvls/eventroad/example/eventpb"
	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/proto"
)

func TestBankAccount_Apply(t *testing.T) {
	testCases := map[string]struct {
		givenEntity BankAccount
		subject     streaming.Subject
		eventBytes  []byte
		wantEntity  BankAccount
		wantError   error
	}{
		"create": {
			subject:     streaming.Subject{Action: "Created"},
			givenEntity: BankAccount{},
			eventBytes:  marshal(&eventpb.Created{Name: "foo"}),
			wantEntity:  BankAccount{Name: "foo"},
		},
		"create uses an empty entity before applying": {
			subject:     streaming.Subject{Action: "Created"},
			givenEntity: BankAccount{Name: "non-fresh-entity"},
			eventBytes:  marshal(&eventpb.Created{Name: ""}),
			wantEntity:  BankAccount{Name: ""},
		},
		"edited": {
			subject:     streaming.Subject{Action: "Edited"},
			givenEntity: BankAccount{Name: "foo"},
			eventBytes:  marshal(&eventpb.Edited{Name: "bar"}),
			wantEntity:  BankAccount{Name: "bar"},
		},
		"edited doesn't override with zero values": {
			subject:     streaming.Subject{Action: "Edited"},
			givenEntity: BankAccount{Name: "foo"},
			eventBytes:  marshal(&eventpb.Edited{Name: ""}),
			wantEntity:  BankAccount{Name: "foo"},
		},
		"unknown event": {
			subject:     streaming.Subject{Action: "Unknown"},
			givenEntity: BankAccount{},
			eventBytes:  marshal(&eventpb.UnhandledEvent{}),
		},
	}
	for testName, tt := range testCases {
		t.Run(testName, func(t *testing.T) {
			got := tt.givenEntity

			err := got.Apply(tt.subject, tt.eventBytes)
			if tt.wantError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantError.Error())
				return
			}
			assert.Equal(t, tt.wantEntity, got)
		})
	}
}

func marshal(pb proto.Message) []byte {
	bytes, err := proto.Marshal(pb)
	if err != nil {
		panic(fmt.Errorf("marshalling %#v: %v", pb, err))
	}
	return bytes
}
