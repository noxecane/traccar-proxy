package config

import (
	"crypto/tls"
	"fmt"

	"github.com/go-pg/pg/v9"
)

func SetupDB(env Env) (*pg.DB, error) {
	opts := &pg.Options{
		Addr:            fmt.Sprintf("%s:%d", env.PostgresHost, env.PostgresPort),
		User:            env.PostgresUser,
		Password:        env.PostgresPassword,
		Database:        env.PostgresDatabase,
		ApplicationName: env.Name,
	}

	if env.PostgresSecureMode {
		opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	db := pg.Connect(opts)
	_, err := db.Exec("select version()")

	return db, err
}
