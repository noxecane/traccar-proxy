package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/tsaron/anansi"
	"tsaron.com/traccar-proxy/pkg/model"
	"tsaron.com/traccar-proxy/pkg/traccar"
)

type latestPositionQuery struct {
	Device uint `key:"device"`
}

type positionQuery struct {
	Device uint      `key:"device"`
	Limit  uint      `key:"limit"`
	Offset uint      `key:"offset"`
	From   time.Time `key:"from"`
	To     time.Time `key:"to"`
}

func Positions(r *chi.Mux, sessions *anansi.SessionStore, repo *traccar.Repo) {
	r.Route("/positions", func(r chi.Router) {
		r.With(sessions.Headless()).Get("/", getPositions(repo))
		r.With(sessions.Headless()).Get("/latest", getLatestPosition(repo))
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

		if !q.From.IsZero() && q.To.IsZero() {
			q.To = time.Now()
		}

		var tps []traccar.Position
		var err error

		if q.From.IsZero() {
			tps, err = repo.FindPositions(r.Context(), q.Device, q.Offset, q.Limit)
		} else {
			tps, err = repo.FindPositionsBetween(r.Context(), q.Device, q.Offset, q.Limit, q.From.UTC(), q.To.UTC())
		}

		if err != nil {
			panic(errors.Wrap(err, "could not get positions"))
		}

		var ps []model.Position
		for _, tp := range tps {
			p, err := traccar.TransformPosition(repo.ToTraccarPosition(&tp))
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
			panic(errors.Wrap(err, "could not get latest position"))
		}

		if p == nil {
			anansi.SendSuccess(r, w, p)
			return
		}

		pos, err := traccar.TransformPosition(repo.ToTraccarPosition(p))
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
