package config

type App struct {
	Name       string `json:"name"`
	Key        string `json:"key"`
	Port       string `json:"port"`
	Host       string `json:"host"`
	Version    string `json:"version"`
	Connection string `json:"connection"`
}
