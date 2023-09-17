package config

type Config struct {
	Port int
}

func Load() (Config, error) {
	cfg := Config{
		Port: 8082,
	}
	return cfg, nil
}
