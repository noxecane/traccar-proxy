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

type Position struct {
	tableName  struct{} `pg:"tc_positions"`
	ID         uint
	CreatedAt  time.Time `pg:"servertime"`
	RecordedAt time.Time `pg:"devicetime"`
	Valid      bool      `pg:",use_zero"`
	Device     uint      `pg:"deviceid"`
	Latitude   float64   `pg:",use_zero"`
	Longitude  float64   `pg:",use_zero"`
	Altitude   float64   `pg:",use_zero"`
	Speed      float64   `pg:",use_zero"`
	Course     float64   `pg:",use_zero"`
	Payload    string    `pg:"attributes"`
	Accuracy   uint
	Address    string
	Protocol   string
	Network    string
	FixedAt    time.Time `pg:"fixtime"`
}

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

func (r *Repo) Device(ctx context.Context, id uint) (*model.Device, error) {
	device := &model.Device{ID: id}
	if err := r.db.ModelContext(ctx, device).WherePK().Select(); err != nil {
		return nil, err
	}

	return device, nil
}

func (r *Repo) Position(ctx context.Context, id uint) (*Position, error) {
	position := &Position{ID: id}
	if err := r.db.ModelContext(ctx, position).WherePK().Select(); err != nil {
		return nil, err
	}

	return position, nil
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

func (r *Repo) Positions(ctx context.Context, device, offset, limit uint) ([]Position, error) {
	positions := []Position{}

	var err error
	if limit == 0 {
		err = r.db.
			ModelContext(ctx, &positions).
			Where("deviceid = ?", device).
			Offset(int(offset)).
			Order("devicetime DESC").
			Select()
	} else {
		err = r.db.
			ModelContext(ctx, &positions).
			Where("deviceid = ?", device).
			Offset(int(offset)).
			Limit(int(limit)).
			Order("devicetime DESC").
			Select()
	}

	return positions, err
}

func (r *Repo) PositionsBetween(ctx context.Context, d, o, l uint, f, t time.Time) ([]Position, error) {
	positions := []Position{}

	var err error
	if l == 0 {
		err = r.db.
			ModelContext(ctx, &positions).
			Where("devicetime [?,?]::tsrange", f.UTC().Format(tsFormat), t.UTC().Format(tsFormat)).
			Where("deviceid = ?", d).
			Offset(int(o)).
			Order("devicetime DESC").
			Select()
	} else {
		err = r.db.
			ModelContext(ctx, &positions).
			Where("devicetime [?,?]::tsrange", f.UTC().Format(tsFormat), t.UTC().Format(tsFormat)).
			Where("deviceid = ?", d).
			Offset(int(o)).
			Limit(int(l)).
			Order("devicetime DESC").
			Select()
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
