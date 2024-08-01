package config

const (
	SERVER = 3000
)

type Config struct {
	Server int `json:"server"`
}

func GetConfig() (Config, error) {
	return Config{
		Server: SERVER,
	}, nil
}
