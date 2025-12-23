package config

type Config struct {
	HTTPAddr  string
	DBPath    string
	JWTSecret string
}

func Load() Config {
	return Config{
		HTTPAddr:  ":8080",
		DBPath:    "C:/DC/mangahub/mangahub.db",
		JWTSecret: "super-secret-key",
	}
}
