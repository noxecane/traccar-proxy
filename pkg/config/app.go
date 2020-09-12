package config

import (
	"net/http"

	"github.com/go-pg/pg/v9"
)

func HealthChecker(db *pg.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		if _, err := db.Exec("select version()"); err != nil {
			http.Error(w, "Could not reach postgres", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		// we don't have a plan for when writes fail
		_, _ = w.Write([]byte("Up and Running!"))
	}
}
