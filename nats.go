package streaming

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	stan "github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

// A NatsServer allows you to use NATS as a `streaming.Server`.
type NatsServer struct {
	Conn stan.Conn
	Log  logrus.Logger
}

// RehydrateEntity replays events on top of the given entity in-place.
// `entity` should be a zero value. Since this operation is used for
// producing new events, only one producer should be rehydrating a subject
// at a time.
func (srv *NatsServer) RehydrateEntity(s Subject, entity ApplyableEntity) error {
	// We need to know the current seq number for this subject in order to
	// know when to stop subscribing.
	seqchan, timeout := make(chan uint64, 1), time.After(1*time.Second)
	subs, err := srv.Conn.Subscribe(SubjectNats(s), func(msg *stan.Msg) {
		seqchan <- msg.Sequence
	}, stan.StartWithLastReceived())
	defer func() { _ = subs.Unsubscribe() }()
	if err != nil {
		return err
	}

	var endSeq uint64
	select {
	case <-timeout:
		return fmt.Errorf("timeout for subject %v", SubjectNats(s))
	case endSeq = <-seqchan:
	}

	// Replay events from the beginning up to endSeq.
	done, timeout := make(chan error, 1), time.After(60*time.Second)
	subs, err = srv.Conn.Subscribe(SubjectNats(s), func(msg *stan.Msg) {
		err = entity.Apply(s, msg.Data)
		if err != nil {
			done <- err
			return
		}
		if msg.Sequence == endSeq {
			done <- nil
			return
		}
	}, stan.DeliverAllAvailable())
	defer func() { _ = subs.Unsubscribe() }()
	if err != nil {
		return err
	}

	select {
	case <-timeout:
		return fmt.Errorf("timeout when looping on subject '%v' until seq '%d'", SubjectNats(s), endSeq)
	case err = <-done:
		if err != nil {
			return err
		}
	}

	return nil
}

// PublishEvent returns errors related to connection and proto.Marshal
// errors.
func (srv *NatsServer) PublishEvent(s Subject, event proto.Message) error {
	bytes, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	err = srv.Conn.Publish(SubjectNats(s), bytes)
	if err != nil {
		return err
	}

	return nil
}

// SubjectNats returns a subject formatted as
// `BankAccount.OTE3Yzk3Y2Y.Created`. If ID or Event are left blank, its
// value is replaced by `*` (wildcard). The Root field cannot be empty.
func SubjectNats(s Subject) string {
	if s.EntityID == "" {
		s.EntityID = "*"
	}
	if s.Action == "" {
		s.Action = "*"
	}
	return fmt.Sprintf("%s.%s.%s", s.Root, s.EntityID, s.Action)
}
