package utl

type DCDN struct {
	DataCenter string
	DomainName string `yaml:"domainName"`
}

type HostAuth struct {
	Name  string
	Uname string `yaml:"host.uname"`
}

type BastionAuth struct {
	Name   string
	Uname  string `yaml:"bastion.uname"`
	Prefix string `yaml:"bastion.prefix"`
}

type SShCfg struct {
	Path string `yaml:"ssh.path"`
	PK   string `yaml:"ssh.pkname"`
	Port int8   `yaml:"ssh.port"`
	Jump bool   `yaml:"ssh.jump"`
}

type HJSShConfig struct {
	HostAuth    *HostAuth
	BastionAuth *BastionAuth
	SSHConfig   *SShCfg
	Dx          *DCDN
}
