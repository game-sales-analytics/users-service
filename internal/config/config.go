package config

type ServerConfig struct {
	Port uint
	Host string
}

type DatabaseConfig struct {
	Port     uint
	Host     string
	Name     string
	Username string
	Password string
	UseAuth  bool
}

type JwtConfig struct {
	Secret string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Jwt      JwtConfig
}
