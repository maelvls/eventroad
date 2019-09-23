package streaming_test

import (
	"fmt"
	"testing"

	natsServer "github.com/nats-io/nats-server/v2/server"
	stanServer "github.com/nats-io/nats-streaming-server/server"
	stan "github.com/nats-io/stan.go"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	runNatsServer(12345)
	_, err := stan.Connect("cluster_id", "client_id", stan.NatsURL(":12345"))
	assert.NoError(t, err, "couldn't connect")

	// Now you can run your tests.
}

// In-memory NATS.
func runNatsServer(port int) *stanServer.StanServer {
	stanOpts := stanServer.GetDefaultOptions()
	stanOpts.ID = "cluster_id"
	s, err := stanServer.RunServerWithOpts(
		stanOpts,
		&natsServer.Options{Port: port},
	)
	if err != nil {
		panic(fmt.Sprintf("cannot launch embedded NATS server: %v", err))
	}
	return s
}
