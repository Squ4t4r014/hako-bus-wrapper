package infrastructure

import (
	"os"
)

type Config struct {
	AbsolutePath string
}

func NewConfig() (*Config, error) {
	conf := new(Config)
	var err error
	conf.AbsolutePath, err = os.Getwd()
	if err != nil {
		return nil, err
	}
	return conf, nil
}
