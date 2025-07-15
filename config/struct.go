package config

// Config 主配置
type Config struct {
	Proxy     GRPCProxyConfig `toml:"proxy"`
	Session   NodeConfig      `toml:"session"`
	Fileapi   NodeConfig      `toml:"fileapi"`
	User      NodeConfig      `toml:"user"`
	Group     NodeConfig      `toml:"group"`
	MsapRead  NodeConfig      `toml:"msap_read"`
	MsapWrite NodeConfig      `toml:"msap_write"`
}

// GRPCProxyConfig grpc Server配置
type GRPCProxyConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
	Log  bool   `toml:"log"`
}

// NodeConfig grpc 远程 配置
type NodeConfig struct {
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	ConnNum int    `toml:"conn_num"`
}
