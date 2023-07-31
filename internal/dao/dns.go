package dao

import "ddns-server/internal/config"

var Dns = dDns{}

type dDns struct{}

func (d *dDns) Get() (*config.Dns, error) {
	cfg, err := config.ReadFromFile(gConfigPath)
	if err != nil {
		return nil, err
	}
	return cfg.Dns, nil
}

func (d *dDns) Set(c *config.Dns) error {
	cfg, err := config.ReadFromFile(gConfigPath)
	if err != nil {
		return err
	}
	cfg.Dns = c
	return config.WriteToFile(cfg, gConfigPath)
}
