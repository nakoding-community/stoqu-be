package config

type Database struct {
	DBHost             string `json:"db_host"`
	DBUser             string `json:"db_user"`
	DBPass             string `json:"db_pass"`
	DBPort             string `json:"db_port"`
	DBName             string `json:"db_name"`
	DBProvider         string `json:"db_provider"`
	DBSSL              string `json:"db_ssl"`
	DBTZ               string `json:"db_tz"`
	AutoMigrateEnabled bool   `json:"automigrate_enabled"`
	SeederEnabled      bool   `json:"seeder_enabled"`
}
