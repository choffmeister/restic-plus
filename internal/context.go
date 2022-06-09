package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

const resticVersion = "0.13.1"

type Context struct {
	ResticVersion string
	ResticBinary  string
	Config        Config
	IdentityFile  string
}

func NewContext(file string) (*Context, error) {
	tempDir := path.Join(os.TempDir(), "restic-plus")
	if file == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		file = path.Join(cwd, "restic-plus.yaml")
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}
	config := Config{}
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("unable to unmarshall config: %w", err)
	}

	resticBinary := path.Join(tempDir, "restic")
	if err := prepareResticBinary(resticVersion, resticBinary); err != nil {
		return nil, fmt.Errorf("unable to prepare restic binary: %w", err)
	}

	identityFile := path.Join(tempDir, "sftp-identity.key")
	if err := prepareIdentity(config.SFTP.IdentityPrivateKey, identityFile); err != nil {
		return nil, fmt.Errorf("unable to prepare identity file: %w", err)
	}

	context := Context{
		ResticVersion: resticVersion,
		ResticBinary:  resticBinary,
		Config:        config,
		IdentityFile:  identityFile,
	}
	return &context, nil
}

func prepareIdentity(identityPrivateKey string, target string) error {
	err := os.MkdirAll(resticDir, 0o755)
	if err != nil {
		return err
	}
	ioutil.WriteFile(sftpIdentityFile, []byte(identityPrivateKey), 0o600)

	return nil
}
