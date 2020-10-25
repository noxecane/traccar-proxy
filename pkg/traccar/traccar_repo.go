package traccar

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/rs/zerolog"
	"tsaron.com/traccar-proxy/pkg/model"
)

const tsFormat = "2006-01-02 15:04"

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

func (r *Repo) FindPositions(ctx context.Context, device, offset, limit uint) ([]Position, error) {
	positions := []Position{}

	query := r.db.
		ModelContext(ctx, &positions).
		Where("deviceid = ?", device).
		Offset(int(offset)).
		Order("devicetime DESC")

	var err error
	if limit == 0 {
		err = query.Select()
	} else {
		err = query.Limit(int(limit)).Select()
	}

	return positions, err
}

func (r *Repo) FindPositionsBetween(ctx context.Context, d, o, l uint, f, t time.Time) ([]Position, error) {
	positions := []Position{}

	query := r.db.
		ModelContext(ctx, &positions).
		Where("devicetime [?,?]::tsrange", f.UTC().Format(tsFormat), t.UTC().Format(tsFormat)).
		Where("deviceid = ?", d).
		Offset(int(o)).
		Order("devicetime DESC")

	var err error
	if l == 0 {
		err = query.Select()
	} else {
		err = query.Limit(int(l)).Select()
	}

	return positions, err
}

func (r *Repo) ToTraccarPosition(p *Position) model.TraccarPosition {
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
