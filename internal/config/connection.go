package config

type Connection struct {
	Primary string `json:"primary"`
	Replica string `json:"replica"`
}
