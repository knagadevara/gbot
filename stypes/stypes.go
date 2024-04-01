package stypes

type HostStr struct {
	HostName   string
	DomainName string
	DataCenter string
	JumpBox    string
}

type SShCfg struct {
	UserName string `yaml:"user_name"`
	SshBase  string `yaml:"ssh_path"`
	KeyName  string `yaml:"ssh_key_name"`
	Jump     bool   `yaml:"jumpstation"`
}
