package dao

import (
	"ddns-server/internal/config"
	"errors"
)

var Record = dRecord{}

type dRecord struct {
}

type Records = map[string]*config.Record

func (d *dRecord) List() ([]*config.Record, error) {
	cfg, err := config.ReadFromFile(gConfigPath)
	if err != nil {
		return nil, err
	}
	return cfg.Records, nil
}

func (d *dRecord) FindByDomain(domain string) (*config.Record, error) {
	rs, err := d.List()
	if err != nil {
		return nil, err
	}
	for _, r := range rs {
		if r.Domain == domain {
			return r, nil
		}
	}

	return nil, errors.New("record not found")
}

func (d *dRecord) Set(record *config.Record, forceUpdateV4 bool, forceUpdateV6 bool) error {
	cfg, err := config.ReadFromFile(gConfigPath)
	if err != nil {
		return err
	}

	update := false
	for i, r := range cfg.Records {
		if r.Domain != record.Domain {
			continue
		}

		update = true
		if record.Ipv4 != "" || forceUpdateV4 {
			cfg.Records[i].Ipv4 = record.Ipv4
		}
		if record.Ipv6 != "" || forceUpdateV6 {
			cfg.Records[i].Ipv6 = record.Ipv6
		}
		break
	}

	if !update {
		cfg.Records = append(cfg.Records, record)
	}

	return config.WriteToFile(cfg, gConfigPath)
}

func (d *dRecord) Delete(domain string) error {
	cfg, err := config.ReadFromFile(gConfigPath)
	if err != nil {
		return err
	}
	newRs := make([]*config.Record, 0)
	for _, r := range cfg.Records {
		if r.Domain == domain {
			continue
		}
		newRs = append(newRs, r)
	}
	cfg.Records = newRs
	return config.WriteToFile(cfg, gConfigPath)
}
