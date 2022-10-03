package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
)

type Configuration struct {
	App        App        `json:"app"`
	Swagger    Swagger    `json:"swagger"`
	JWT        JWT        `json:"jwt"`
	Databases  []Database `json:"databases"`
	Connection Connection `json:"connection"`
	Driver     Driver     `json:"driver"`
}

var Config *Configuration = &Configuration{}

func Load(path string) (*Configuration, error) {
	if path == "" {
		wd, err := os.Getwd()
		if err != nil {
			return Config, err
		}
		path = fmt.Sprintf("%s/config/config.%s.yml", wd, os.Getenv("ENV"))
	}

	if os.Getenv("ENV") != "development" && os.Getenv("ENV") != "debug" && os.Getenv("ENV") != "local" {
		path = fmt.Sprintf("/run/secrets/%s", os.Getenv("CONFIG"))
	}

	err := configor.New(&configor.Config{AutoReload: true, AutoReloadInterval: time.Minute}).Load(Config, path)
	if err != nil {
		logrus.Info(err)
		return Config, err
	}

	return Config, nil
}

func (Configuration) String() string {
	sb := strings.Builder{}
	return sb.String()
}

func (c *Configuration) Raw() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
