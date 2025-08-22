package config

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Static   StaticConfig   `mapstructure:"static"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type StaticConfig struct {
	RootFolder   string `mapstructure:"root_folder"`
	BackupFolder string `mapstructure:"backup_folder"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}
