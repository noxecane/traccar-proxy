package proxy

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"tsaron.com/positions/pkg/model"
	"tsaron.com/positions/pkg/traccar"
)

type Emitter struct {
	log  zerolog.Logger
	repo *traccar.Repo
	conn *nats.EncodedConn
}

type PositionEvent struct {
	Action   string                `json:"action"`
	Position model.TraccarPosition `json:"data"`
}

// TODO: https://github.com/nats-io/nats.go/blob/master/examples/nats-qsub/main.go (options)

func NewEmitter(conn *nats.Conn, repo *traccar.Repo, log zerolog.Logger) (*Emitter, error) {
	subLogger := log.With().Str("source", "emitter").Logger()
	encConn, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	return &Emitter{subLogger, repo, encConn}, nil
}

func (e *Emitter) Run(ctx context.Context, wg *sync.WaitGroup) {
	// WaitGroup to force blocking on the caller
	wg.Add(1)

	out := make(chan []byte, 64)

	// start listening for events
	go e.repo.Listen(ctx, "tc_positions", out)

	go func() {
		for {
			select {
			case ev := <-out:
				var event PositionEvent
				if err := json.Unmarshal(ev, &event); err != nil {
					e.log.Err(err).RawJSON("event", ev).Msg("failed to to decode event")
					continue
				}

				if event.Action != "INSERT" {
					continue
				}

				p := event.Position

				res, err := traccar.TransformPosition(p)
				if err != nil {
					e.log.Err(err).Interface("position", p).Msg("")
					continue
				}

				if err := e.conn.Publish("traccar.positions", res); err != nil {
					e.log.Err(err).Interface("position", res).Msg("failed to publish")
				}

			case <-ctx.Done():
				e.log.Info().Msg("shutting down the emitter")

				// draining the nats connection
				if err := e.conn.Drain(); err != nil {
					e.log.Err(err).Msg("failed to drain nats connection")
				}

				// tell the caller we're done
				wg.Done()
			}
		}
	}()
}
