package config

import "github.com/nats-io/nats.go"

func SetupNats(env Env) (*nats.Conn, error) {
	opts := []nats.Option{nats.Name(env.Name)}

	if env.NatsUser != "" {
		opts = append(opts, nats.UserInfo(env.NatsUser, env.NatsPassword))
	}

	nc, err := nats.Connect(env.NatsUrl, opts...)

	return nc, err
}
