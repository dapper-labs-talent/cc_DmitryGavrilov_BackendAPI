package config

func CreateTestConfiguration() *Config {
	return &Config{
		Server: Server{
			ListenPort: 8080,
		},
		JWT: JWT{
			Expiration: 10,
			Secret:     "Test",
		},
		Database: Database{
			Driver: "memory",
		},
	}
}
