package config

import (
	"fmt"
)

// Get information function
// - GetServerAddress: get server address
// - GetStaticFolder: get static folder
// - GetBackupFolder: get backup folder
// - GetHost: get host
// - GetPort: get port
// - GerDatabaseInformation: get db host, db port, db name
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf(":%s", c.Server.Port)
}

func (c *Config) GetStaticFolder() string {
	return c.Static.RootFolder
}

func (c *Config) GetBackupFolder() string {
	return c.Static.BackupFolder
}

func (c *Config) GetHost() string {
	return c.Server.Host
}

func (c *Config) GetPort() string {
	return c.Server.Port
}

func (c *Config) GerDatabaseInformation() string {
	return fmt.Sprintf("name:%s, host:%s, port:%s", c.Database.Database, c.Database.Host, c.Database.Port)
}

// show config
func PrintConfig(cfg *Config, env string) {
	fmt.Println("=== Current Configuration ===")
	fmt.Printf("environment: %s\n", env)
	fmt.Printf("server address: %s\n", cfg.GetServerAddress())
	fmt.Printf("static folder: %s\n", cfg.GetStaticFolder())
	fmt.Printf("backup folder: %s\n", cfg.GetBackupFolder())
	fmt.Printf("database information: %s\n", cfg.GerDatabaseInformation())
	fmt.Println("================")
}
