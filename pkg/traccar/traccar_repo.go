package traccar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/rs/zerolog"
	"tsaron.com/traccar-proxy/pkg/model"
)

const pgTimef = "2006-01-02 15:04"

var ErrInvalidQuery = errors.New("your query is invalid")

type Repo struct {
	log     zerolog.Logger
	db      *pg.DB
	channel string
}

func NewRepo(db *pg.DB, eventChannel string, log zerolog.Logger) *Repo {
	subLogger := log.With().Str("source", "traccar-repo").Logger()
	return &Repo{subLogger, db, eventChannel}
}

type tableEvent struct {
	Table  string      `json:"table"`
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

func (r *Repo) Listen(ctx context.Context, table string, out chan<- []byte) {
	l := r.db.Listen(r.channel)
	defer l.Close()

	ch := l.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case n := <-ch:
			e := new(tableEvent)

			if err := json.Unmarshal([]byte(n.Payload), e); err != nil {
				r.log.Err(err).Msg("failed to decode event payload")
				continue
			}

			if e.Table != table {
				continue
			}

			r.log.
				Info().
				Str("table", e.Table).
				Str("action", e.Action).
				Interface("data", e.Data).
				Msg("received event")

			out <- []byte(n.Payload)
		}
	}
}

func (r *Repo) FindDevice(ctx context.Context, externalID string) (*Device, error) {
	device := new(Device)

	err := r.db.
		ModelContext(ctx, device).
		Where("uniqueid = ?", externalID).
		Select()

	if err == pg.ErrNoRows {
		return nil, nil
	}

	return device, err
}

func (r *Repo) LatestPosition(ctx context.Context, device uint) (*Position, error) {
	position := &Position{}
	err := r.db.
		ModelContext(ctx, position).
		Where("deviceid = ?", device).
		Order("devicetime DESC").
		Limit(1).
		Select()

	if err == pg.ErrNoRows {
		return nil, nil
	}

	return position, err
}

type QueryOpts struct {
	// The oldest position by devicetime. If this is not set, time based query is entirely ignroed
	From time.Time
	// The latest position by devicetime.
	To time.Time
	// The maximum number of positions to return
	Limit int
	// How many positions to skip before limit starts getting counted
	Offset int
	// Order of the results. Could be "oldest" meaning oldest first or "latest" meaning latest first
	// arranged by the devicetime of the positions.
	Order string
}

func (r *Repo) FindPositions(ctx context.Context, device uint, opts QueryOpts) ([]Position, error) {
	positions := []Position{}

	order := "devicetime"
	switch opts.Order {
	case "oldest":
		order += " ASC"
	case "latest":
		order += " DESC"
	default:
		return nil, ErrInvalidQuery
	}

	query := r.db.
		ModelContext(ctx, &positions).
		Where("deviceid = ?", device).
		Offset(opts.Offset).
		Order(order)

	switch {
	case !opts.From.IsZero() && !opts.To.IsZero():
		tRange := fmt.Sprintf("[%s, %s]", opts.From.Format(pgTimef), opts.To.Format(pgTimef))
		query = query.Where("?::tsrange @> devicetime", tRange)
	case opts.From.IsZero() && opts.To.IsZero():
		// no-op
	default:
		return nil, ErrInvalidQuery
	}

	var err error
	if opts.Limit == 0 {
		err = query.Select()
	} else {
		err = query.Limit(opts.Limit).Select()
	}

	return positions, err
}

func (r *Repo) RemoveTZ(p *Position) model.TraccarPosition {
	return model.TraccarPosition{
		ID:         p.ID,
		CreatedAt:  model.ISOWithoutTZ(p.CreatedAt),
		RecordedAt: model.ISOWithoutTZ(p.RecordedAt),
		Valid:      p.Valid,
		Device:     p.Device,
		Latitude:   p.Latitude,
		Longitude:  p.Longitude,
		Altitude:   p.Altitude,
		Speed:      p.Speed,
		Course:     p.Course,
		Payload:    p.Payload,
		Accuracy:   p.Accuracy,
		Address:    p.Address,
		Protocol:   p.Protocol,
		Network:    p.Network,
		FixedAt:    model.ISOWithoutTZ(p.FixedAt),
	}
}
