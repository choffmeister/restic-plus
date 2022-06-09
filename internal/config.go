package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Targets []string     `yaml:"targets"`
	Restic  ConfigRestic `yaml:"restic"`
	SFTP    ConfigSFTP   `yaml:"sftp"`
	Cron    ConfigCron   `yaml:"cron"`
}

type ConfigRestic struct {
	Password string `yaml:"password"`
}

type ConfigSFTP struct {
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	User               string `yaml:"user"`
	IdentityPrivateKey string `yaml:"identityPrivateKey"`
}

type ConfigCron struct {
	Cleanup ConfigCronCleanup `yaml:"cleanup"`
}

type ConfigCronCleanup struct {
	Enabled bool                  `yaml:"enabled"`
	Keep    ConfigCronCleanupKeep `yaml:"keep"`
}

type ConfigCronCleanupKeep struct {
	Last    int `yaml:"last"`
	Daily   int `yaml:"daily"`
	Weekly  int `yaml:"weekly"`
	Monthly int `yaml:"monthly"`
	Yearly  int `yaml:"yearly"`
}

type ConfigBandwidth struct {
	Download int `yaml:"download"`
	Upload   int `yaml:"upload"`
}

func (c *Config) LoadFromFile(dir string) error {
	if dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		dir = cwd
	}

	file := path.Join(dir, "restic-plus.yaml")
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("unable to read config: %w", err)
	}

	if err := yaml.Unmarshal(bytes, c); err != nil {
		return fmt.Errorf("unable to unmarshall config: %w", err)
	}

	return nil
}
