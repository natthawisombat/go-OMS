package pkglogger

var L = new(config)

type config struct {
	LogPath    string `env:"LOG_PATH,default=./assets" json:",omitempty"`
	LogLevel   string `env:"LOG_LEVEL" json:",omitempty"`
	LogMode    string `env:"LOG_MODE,default=development" json:",omitempty"`
	LogSize    int    `env:"LOG_SIZE,default=10" json:",omitempty"`
	LogBackups int    `env:"LOG_BACKUPS,default=3" json:",omitempty"`
	LogAge     int    `env:"LOG_AGE,default=45" json:",omitempty"`
}
