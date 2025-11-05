package config

type serverConfig struct {
	Port               string
	Role               string
	ReplicaOf          string
	Master_replid      string
	Master_repl_offset int64
}

var ServerConfig *serverConfig

func InitServerConfig() {
	ServerConfig = &serverConfig{
		Port:               "6379",
		Role:               "master",
		ReplicaOf:          "",
		Master_replid:      "",
		Master_repl_offset: 0,
	}
}
