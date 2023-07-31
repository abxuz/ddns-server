package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/xbugio/b-tools/bslice"
	"gopkg.in/yaml.v3"
)

var gLock sync.RWMutex

type Auth struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type Web struct {
	Listen string `yaml:"listen" json:"listen"`
	Auth   *Auth  `yaml:"auth,omitempty" json:"auth,omitempty"`
}

type Dns struct {
	Listen string `yaml:"listen" json:"listen"`
}

type Record struct {
	Domain string `yaml:"domain" json:"domain"`
	Ipv4   string `yaml:"ipv4" json:"ipv4"`
	Ipv6   string `yaml:"ipv6" json:"ipv6"`
}

type Config struct {
	Dns     *Dns      `yaml:"dns,omitempty"`
	Web     *Web      `yaml:"web,omitempty"`
	Records []*Record `yaml:"records"`
}

func ReadFromFile(p string) (*Config, error) {
	gLock.RLock()
	defer gLock.RUnlock()

	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	c := new(Config)
	err = yaml.NewDecoder(f).Decode(c)
	f.Close()
	if err != nil {
		return nil, err
	}

	if err := c.Valid(); err != nil {
		return nil, err
	}
	return c, nil
}

func WriteToFile(c *Config, p string) error {
	if err := c.Valid(); err != nil {
		return err
	}

	gLock.Lock()
	defer gLock.Unlock()

	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	defer encoder.Close()
	return encoder.Encode(c)
}

func (c *Config) Valid() error {
	if c.Dns != nil {
		if err := c.Dns.Valid(); err != nil {
			return fmt.Errorf("err in dns config: %v", err.Error())
		}
	}

	if c.Web != nil {
		if err := c.Web.Valid(); err != nil {
			return fmt.Errorf("err in web config: %v", err.Error())
		}
	}

	ok := bslice.Unique(len(c.Records), func(i int) any { return c.Records[i].Domain })
	if !ok {
		return fmt.Errorf("duplicate domain name")
	}
	return nil
}

func (c *Dns) Valid() error {
	if c.Listen == "" {
		return fmt.Errorf("missing listen")
	}
	return nil
}

func (c *Web) Valid() error {
	if c.Listen == "" {
		return fmt.Errorf("missing listen")
	}
	return nil
}
