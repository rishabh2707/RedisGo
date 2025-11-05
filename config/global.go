package config

type serverConfig struct {
	Role               string
	ReplicaOf          string
	Master_replid      string
	Master_repl_offset int64
}

var ServerConfig *serverConfig

func InitServerConfig() {
	ServerConfig = &serverConfig{
		Role:               "master",
		ReplicaOf:          "",
		Master_replid:      "",
		Master_repl_offset: 0,
	}
}
