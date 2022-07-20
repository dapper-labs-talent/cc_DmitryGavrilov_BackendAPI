package config

type Config struct {
	Server
	Postgres
	JWT
	Logging
}

type JWT struct {
}

type Postgres struct {
}

type Server struct {
}

type Logging struct {
}
