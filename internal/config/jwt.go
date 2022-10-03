package config

type JWT struct {
	Secret        string `json:"secret"`
	RefreshSecret string `json:"refresh_secret"`
}
