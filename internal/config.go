package internal

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Stanzas   []ConfigStanza  `yaml:"stanzas"`
	Restic    ConfigRestic    `yaml:"restic"`
	SFTP      ConfigSFTP      `yaml:"sftp"`
	Cron      ConfigCron      `yaml:"cron"`
	Bandwidth ConfigBandwidth `yaml:"bandwidth"`
}

type ConfigStanza struct {
	Type           string `yaml:"type"`
	Implementation Stanza
}

func (ct *ConfigStanza) UnmarshalYAML(value *yaml.Node) error {
	type rawConfigStanza ConfigStanza
	if err := value.Decode((*rawConfigStanza)(ct)); err != nil {
		return err
	}

	switch ct.Type {
	case "":
		fallthrough
	case DirectoryStanzaType:
		implementation := &DirectoryStanza{}
		if err := value.Decode(implementation); err != nil {
			return fmt.Errorf("invalid configuration for stanza %s: %w", ct.Type, err)
		}
		ct.Implementation = implementation
	case ZFSDatasetStanzaType:
		implementation := &ZFSDatasetStanza{}
		if err := value.Decode(implementation); err != nil {
			return fmt.Errorf("invalid configuration for stanza %s: %w", ct.Type, err)
		}
		ct.Implementation = implementation
	case ZFSZvolStanzaType:
		implementation := &ZFSZvolStanza{}
		if err := value.Decode(implementation); err != nil {
			return fmt.Errorf("invalid configuration for stanza %s: %w", ct.Type, err)
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
