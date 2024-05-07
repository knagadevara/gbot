package utl

type SShCfg struct {
	Path string `yaml:"ssh.path"`
	PK   string `yaml:"ssh.pkname"`
	Port int8   `yaml:"ssh.port"`
	Jump bool   `yaml:"ssh.jump"`
}

type HostAuth struct {
	Name  string `yaml:"hostname"`
	Uname string `yaml:"host.uname"`
}

type BastionAuth struct {
	Name   string `yaml:"bastion.hostname"`
	Uname  string `yaml:"bastion.uname"`
	Prefix string `yaml:"bastion.prefix"`
}

type DCDN struct {
	DataCenter string `yaml:"datacenter"`
	DomainName string `yaml:"domainName"`
}

type HJSShConfig struct {
	HostAuth    *HostAuth    `yaml:"host"`
	BastionAuth *BastionAuth `yaml:"bastion"`
	SSHConfig   *SShCfg      `yaml:"ssh"`
	Dx          *DCDN        `yaml:"dcdn"`
}
