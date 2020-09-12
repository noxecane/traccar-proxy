package config

// Env is the expected config values from the process's environment
type Env struct {
	AppEnv string `default:"dev" split_words:"true"`
	Name   string `required:"true"`
	Port   int    `required:"true"`
	Scheme string `required:"true"`
	Secret []byte `required:"true"`

	PostgresHost       string `required:"true" split_words:"true"`
	PostgresPort       int    `required:"true" split_words:"true"`
	PostgresSecureMode bool   `required:"true" split_words:"true"`
	PostgresUser       string `required:"true" split_words:"true"`
	PostgresPassword   string `required:"true" split_words:"true"`
	PostgresDatabase   string `required:"true" split_words:"true"`

	HeadlessTimeout string `required:"true" split_words:"true"`
}
