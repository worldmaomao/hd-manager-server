package config

type Configuration struct {
	Service  Service
	Database Database
}

func (config *Configuration) GetJwtKey() string {
	return config.Service.JwtKey
}

type Service struct {
	Host   string
	Port   int
	JwtKey string
}

type Database struct {
	Host         string
	Port         int
	DatabaseName string
	Username     string
	Password     string
}
