package dao

import "ddns-server/internal/config"

var Web = dWeb{}

type dWeb struct{}

func (d *dWeb) Get() (*config.Web, error) {
	cfg, err := config.ReadFromFile(gConfigPath)
	if err != nil {
		return nil, err
	}
	return cfg.Web, nil
}

func (d *dWeb) Set(c *config.Web) error {
	cfg, err := config.ReadFromFile(gConfigPath)
	if err != nil {
		return err
	}
	cfg.Web = c
	return config.WriteToFile(cfg, gConfigPath)
}
