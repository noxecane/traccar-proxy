package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/tsaron/anansi"
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
		r.Get("/latest", getLatestPosition(repo))
	})
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
