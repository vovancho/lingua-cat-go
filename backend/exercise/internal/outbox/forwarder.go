package outbox

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
)

func StartForwarder(ctx context.Context, source message.Subscriber, destination message.Publisher) error {
	logger := watermill.NewStdLogger(false, false)

	fwd, err := forwarder.NewForwarder(
		source,
		destination,
		logger,
		forwarder.Config{
			ForwarderTopic: "lcg.exercise.completed",
		},
	)
	if err != nil {
		return err
	}

	go func() {
		if err := fwd.Run(ctx); err != nil {
			panic(err)
		}
	}()

	return nil
}
