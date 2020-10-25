package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	chiWare "github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v9"
	"github.com/tsaron/anansi"
	"github.com/tsaron/anansi/middleware"
	"tsaron.com/traccar-proxy/pkg/config"
	"tsaron.com/traccar-proxy/pkg/proxy"
	"tsaron.com/traccar-proxy/pkg/rest"
	"tsaron.com/traccar-proxy/pkg/traccar"
)

func main() {
	var err error

	var env config.Env
	if err = anansi.LoadEnv(&env); err != nil {
		panic(err)
	}

	log := anansi.NewLogger(env.Name)

	// connect to postgresql
	var db *pg.DB
	if db, err = config.SetupDB(env); err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Err(err).Msg("failed to disconnect from postgres cleanly")
		}
	}()
	log.Info().Msg("successfully connected to postgres")

	nc, err := config.SetupNats(env)
	if err != nil {
		panic(err)
	}
	log.Info().Msg("successfully connected to nats server")

	repo := traccar.NewRepo(db, "traccar.events", log)

	sessions := anansi.NewSessionStore(env.Secret, env.Scheme, 0, nil)

	// API router
	router := chi.NewRouter()

	// setup app middlware
	middleware.CORS(router, env.AppEnv, "https://*.tsaron.com", "https://*castui.netlify.app", "http://localhost:8080")
	middleware.DefaultMiddleware(router)
	router.Use(middleware.AttachLogger(log))
	router.Use(middleware.TrackRequest())
	router.Use(middleware.TrackResponse())
	router.Use(middleware.Recoverer(env.AppEnv))
	router.Use(chiWare.Timeout(time.Minute))

	router.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Whoops!! This route doesn't exist", http.StatusNotFound)
	})

	rest.Positions(router, sessions, repo)
	rest.Devices(router, sessions, repo)

	// mount API on app router
	appRouter := chi.NewRouter()
	appRouter.Mount("/api/v1/traccar", router)
	appRouter.Get("/", config.HealthChecker(db))
	appRouter.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Whoops!! This route doesn't exist", http.StatusNotFound)
	})

	// run server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	emitter, err := proxy.NewEmitter(nc, repo, log)
	if err != nil {
		panic(err)
	}

	done := new(sync.WaitGroup)

	emitter.Run(ctx, done)

	go anansi.CancelOnInterrupt(cancel, log)
	anansi.RunServer(ctx, log, &http.Server{
		Addr:    fmt.Sprintf(":%d", env.Port),
		Handler: appRouter,
	})

	// make sure the worker dies
	done.Wait()
}
