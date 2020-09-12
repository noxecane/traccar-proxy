package traccar

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/rs/zerolog"
	"tsaron.com/positions/pkg/model"
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

type TableEvent struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

func (r *Repo) Listen(ctx context.Context, table string, out chan<- TableEvent) {
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

			out <- TableEvent{e.Action, e.Data}
		}
	}
}

func (r *Repo) Device(ctx context.Context, id uint) (*model.Device, error) {
	device := &model.Device{ID: id}
	if err := r.db.ModelContext(ctx, device).WherePK().Select(); err != nil {
		return nil, err
	}

	return device, nil
}

func (r *Repo) Position(ctx context.Context, id uint) (*model.Position, error) {
	position := &model.Position{ID: id}
	if err := r.db.ModelContext(ctx, position).WherePK().Select(); err != nil {
		return nil, err
	}

	return position, nil
}

func (r *Repo) LatestPosition(ctx context.Context, device uint) (*model.Position, error) {
	position := &model.Position{}
	err := r.db.
		ModelContext(ctx, position).
		Where("deviceid = ?", device).
		Order("servertime DESC").
		Limit(1).
		Select()

	return position, err
}

func (r *Repo) Positions(ctx context.Context, device uint, limit uint) ([]model.Position, error) {
	positions := []model.Position{}

	var err error
	if limit == 0 {
		err = r.db.
			ModelContext(ctx, &positions).
			Where("deviceid = ?", device).
			Select()
	} else {
		err = r.db.
			ModelContext(ctx, &positions).
			Where("deviceid = ?", device).
			Limit(int(limit)).
			Select()
	}

	return positions, err
}

func (r *Repo) PositionsBetween(ctx context.Context, d uint, f, t time.Time) ([]model.Position, error) {
	positions := []model.Position{}
	err := r.db.
		ModelContext(ctx, &positions).
		Where("servertime [?,?]::tsrange", f.UTC().Format(tsFormat), t.UTC().Format(tsFormat)).
		Where("deviceid = ?", d).
		Select()

	return positions, err
}
