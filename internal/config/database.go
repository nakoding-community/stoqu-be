package config

type Database struct {
	DBHost        string `json:"db_host"`
	DBUser        string `json:"db_user"`
	DBPass        string `json:"db_pass"`
	DBPort        string `json:"db_port"`
	DBName        string `json:"db_name"`
	DBProvider    string `json:"db_provider"`
	DBSSL         string `json:"db_ssl"`
	DBTZ          string `json:"db_tz"`
	DBAutoMigrate bool   `json:"db_automigrate"`
	DBSeeder      bool   `json:"db_seeder"`
	DBLogLevel    int    `json:"db_loglevel"`
}
