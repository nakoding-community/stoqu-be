package config

type Swagger struct {
	SwaggerHost   string `enjsonv:"swagger_host"`
	SwaggerScheme string `enjsonv:"swagger_scheme"`
	SwaggerPrefix string `json:"swagger_prefix"`
}
