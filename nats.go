package streaming

import (
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
// `entity` should be a zero value.
func (srv NatsServer) RehydrateEntity(s Subject, entity ApplyableEntity) error {
	var err error
	handler := func(msg *stan.Msg) {
		// Apply mutates `entity` itself.
		err := entity.Apply(s, msg.Data)
		if err != nil {
			return
		}
	}

	srv.Conn.Subscribe(s.String(), handler, stan.DeliverAllAvailable())

	return err
}

// PublishEvent returns errors related to connection and proto.Marshal
// errors.
func (srv NatsServer) PublishEvent(s Subject, event proto.Message) error {
	bytes, err := proto.Marshal(event)
	if err != nil {
		return err
	}
	srv.Conn.Publish(s.String(), bytes)
}
