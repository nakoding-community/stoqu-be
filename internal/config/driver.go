package config

type (
	Driver struct {
		Firestore     Firestore     `json:"firestore"`
		Elasticsearch Elasticsearch `json:"elasticsearch"`
		Sentry        Sentry        `json:"dsn"`
	}

	Firestore struct {
		Credentials string `json:"credentials"`
		ProjectID   string `json:"project_id"`
	}

	Elasticsearch struct {
		Url      string `json:"url"`
		User     string `json:"user"`
		Password string `json:"password"`
	}
	Sentry struct {
		Dsn string `json:"url"`
	}
)
