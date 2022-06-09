package internal

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Targets   []ConfigTarget  `yaml:"targets"`
	Restic    ConfigRestic    `yaml:"restic"`
	SFTP      ConfigSFTP      `yaml:"sftp"`
	Cron      ConfigCron      `yaml:"cron"`
	Bandwidth ConfigBandwidth `yaml:"bandwidth"`
}

type ConfigTarget struct {
	Type           string `yaml:"type"`
	Implementation Target
}

func (ct *ConfigTarget) UnmarshalYAML(value *yaml.Node) error {
	type rawConfigTarget ConfigTarget
	if err := value.Decode((*rawConfigTarget)(ct)); err != nil {
		return err
	}

	switch ct.Type {
	case "":
		fallthrough
	case DirectoryTargetType:
		implementation := &DirectoryTarget{}
		if err := value.Decode(implementation); err != nil {
			return fmt.Errorf("invalid configuration for target %s: %w", ct.Type, err)
		}
		ct.Implementation = implementation
	case ZFSDatasetTargetType:
		implementation := &ZFSDatasetTarget{}
		if err := value.Decode(implementation); err != nil {
			return fmt.Errorf("invalid configuration for target %s: %w", ct.Type, err)
		}
		ct.Implementation = implementation
	case ZFSZvolTargetType:
		implementation := &ZFSZvolTarget{}
		if err := value.Decode(implementation); err != nil {
			return fmt.Errorf("invalid configuration for target %s: %w", ct.Type, err)
		}
		ct.Implementation = implementation
	default:
		return fmt.Errorf("unknown type %s", ct.Type)
	}

	return nil
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
