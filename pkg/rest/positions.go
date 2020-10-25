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
	Limit  int       `key:"limit"`
	Offset int       `key:"offset"`
	From   time.Time `key:"from"`
	To     time.Time `key:"to"`
	Order  string    `key:"order" default:"latest"`
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

		tps, err := repo.FindPositions(r.Context(), q.Device, traccar.QueryOpts{
			From:   q.From,
			To:     q.To,
			Offset: q.Offset,
			Limit:  q.Limit,
			Order:  q.Order,
		})
		if err != nil {
			panic(errors.Wrap(err, "could not get positions"))
		}

		var ps []model.Position
		for _, tp := range tps {
			p, err := traccar.TransformPosition(repo.RemoveTZ(&tp))
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

		pos, err := traccar.TransformPosition(repo.RemoveTZ(p))
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
