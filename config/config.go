package config

type Config struct {
	Host   string            `json:"host"`
	Port   string            `json:"port"`
	IsAuth bool              `json:"is_auth"`
	User   map[string]string `json:"user"`
}
