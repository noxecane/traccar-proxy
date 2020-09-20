package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tsaron/anansi"
	"tsaron.com/positions/pkg/model"
	"tsaron.com/positions/pkg/traccar"
)

func Devices(r *chi.Mux, sessions *anansi.SessionStore, repo *traccar.Repo) {
	r.Route("/devices", func(r chi.Router) {
		r.With(sessions.Headless()).Get("/{externalID}", getDevice(repo))
	})
}

func getDevice(repo *traccar.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		externalID := anansi.StringParam(r, "externalID")

		dev, err := repo.FindDevice(r.Context(), externalID)
		if err != nil {
			panic(err)
		}

		if dev == nil {
			panic(anansi.APIError{
				Code:    http.StatusNotFound,
				Message: "Could not find device with the given ID",
			})
		}

		anansi.SendSuccess(r, w, model.Device{
			ID:         dev.ID,
			Name:       dev.Name,
			ExternalID: dev.ExternalID,
		})
	}
}
