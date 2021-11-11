package config

func getDefaults() Config {
	return Config{
		Server: ServerConfig{
			Port: 50050,
			Host: "127.0.0.1",
		},
		Database: DatabaseConfig{
			Port:     27018,
			Host:     "users_db",
			Name:     "users",
			Username: "",
			Password: "",
			UseAuth:  false,
		},
		Jwt: JwtConfig{
			Key: "",
		},
	}
}
