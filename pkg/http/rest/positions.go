package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/tsaron/anansi"
	"tsaron.com/positions/pkg/model"
	"tsaron.com/positions/pkg/traccar"
)

type latestPositionQuery struct {
	Device uint `json:"device"`
}

type positionQuery struct {
	Device uint      `json:"device"`
	Limit  uint      `json:"limit"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
}

func Positions(r *chi.Mux, repo *traccar.Repo) {
	r.Route("/positions", func(r chi.Router) {
		r.Get("/", getPositions(repo))
		r.Get("/latest", getLatestPosition(repo))
	})
}

func getPositions(repo *traccar.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := new(positionQuery)
		anansi.ReadQuery(r, q)

		if q.Device == 0 {
			panic(anansi.APIError{
				Code:    http.StatusBadRequest,
				Message: "You need to pass a device ID",
			})
		}

		if q.From.IsZero() != q.To.IsZero() {
			panic(anansi.APIError{
				Code:    http.StatusBadRequest,
				Message: "Both from and to must be set if any is set at all",
			})
		}

		var tps []model.TraccarPosition
		var err error

		if q.From.IsZero() {
			tps, err = repo.Positions(r.Context(), q.Device, q.Limit)
		} else {
			tps, err = repo.PositionsBetween(r.Context(), q.Device, q.From.UTC(), q.To.UTC())
		}

		if err != nil {
			panic(err)
		}

		var ps []model.Position
		for _, tp := range tps {
			p, err := traccar.TransformPosition(tp)
			if err != nil {
				panic(anansi.APIError{
					Code:    http.StatusUnprocessableEntity,
					Message: "Could not parse position because of attribute",
					Err:     err,
					Meta:    p, // send this as it's still useful
				})
			}
			ps = append(ps, p)
		}

		anansi.SendSuccess(r, w, ps)
	}
}

func getLatestPosition(repo *traccar.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := new(latestPositionQuery)
		anansi.ReadQuery(r, q)

		if q.Device == 0 {
			panic(anansi.APIError{
				Code:    http.StatusBadRequest,
				Message: "You need to pass a device ID",
			})
		}

		p, err := repo.LatestPosition(r.Context(), q.Device)
		// let anansi take care of the error
		if err != nil {
			panic(err)
		}

		pos, err := traccar.TransformPosition(*p)
		if err != nil {
			panic(anansi.APIError{
				Code:    http.StatusUnprocessableEntity,
				Message: "Could not parse position because of attribute",
				Err:     err,
				Meta:    pos, // send this as it's still useful
			})
		}

		anansi.SendSuccess(r, w, pos)
	}
}
