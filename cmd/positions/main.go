package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	chiWare "github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v9"
	"github.com/nats-io/nats.go"
	"github.com/tsaron/anansi"
	"github.com/tsaron/anansi/middleware"
	"tsaron.com/positions/pkg/config"
	"tsaron.com/positions/pkg/proxy"
	"tsaron.com/positions/pkg/traccar"
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

	nc, err := nats.Connect(env.NatsUrl)
	if err != nil {
		panic(err)
	}
	log.Info().Msg("successfully connected to nats server")

	// var tmt time.Duration
	// if tmt, err = time.ParseDuration(env.HeadlessTimeout); err != nil {
	// 	panic(err)
	// }
	// sessions := anansi.NewSessionStore(env.Secret, env.Scheme, tmt, nil)

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

	// mount API on app router
	appRouter := chi.NewRouter()
	appRouter.Mount("/api/v1", router)
	appRouter.Get("/", config.HealthChecker(db))
	appRouter.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Whoops!! This route doesn't exist", http.StatusNotFound)
	})

	// run server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go anansi.CancelOnInterrupt(cancel, log)

	repo := traccar.NewRepo(db, "traccar.events", log)

	emitter, err := proxy.NewEmitter(nc, repo, log)
	if err != nil {
		panic(err)
	}

	done := emitter.Run(ctx)

	anansi.RunServer(ctx, log, &http.Server{
		Addr:    fmt.Sprintf(":%d", env.Port),
		Handler: appRouter,
	})

	// make sure the worker dies
	done.Wait()
}
