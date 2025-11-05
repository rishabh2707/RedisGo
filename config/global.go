package config

type serverConfig struct {
	Role      string
	ReplicaOf string
}

var ServerConfig *serverConfig

func InitServerConfig() {
	ServerConfig = &serverConfig{
		Role:      "master",
		ReplicaOf: "",
	}
}
