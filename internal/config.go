package internal

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
