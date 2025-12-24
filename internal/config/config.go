package config

type Config struct {
	HTTPAddr  string
	GRPCAddr  string
	DBPath    string
	JWTSecret string
}

func Load() Config {
	return Config{
		HTTPAddr:  ":8080",
		GRPCAddr:  ":50051",
		DBPath:    "C:/DC/mangahub/mangahub.db",
		JWTSecret: "super-secret-key",
	}
}
