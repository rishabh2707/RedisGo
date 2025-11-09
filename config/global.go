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

type rdbConfig struct {
	Dir        string
	DbFileName string
}

var RDBConfig *rdbConfig

func InitRDBConfig() {
	RDBConfig = &rdbConfig{
		Dir:        "/tmp/redis-files",
		DbFileName: "dump.rdb",
	}
}
