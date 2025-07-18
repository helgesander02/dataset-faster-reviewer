package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Environment type
// - Development: development environment
// - Production: production environment
// - Other: other environment
func LoadConfigForEnvironment(env string) (*Config, error) {
	configName := fmt.Sprintf("config.%s", env)
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/yaml")
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Configuration file %s.yaml not found, use defaults\n", configName)
		} else {
			return nil, fmt.Errorf("Failed to read configuration file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("Failed to parse configuration: %w", err)
	}
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("Configuration verification failed: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// server config default value
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "localhost")

	// static folder default value
	viper.SetDefault("static.folder", "./static")

	// database config default value
	viper.SetDefault("database.host", "")
	viper.SetDefault("database.port", "")
	viper.SetDefault("database.username", "")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.database", "")
}

// Configuration Verification 
func validateConfig(config *Config) error {
	if config.Static.Folder == "" {
		return fmt.Errorf("STATIC_FOLDER is empty")
	}
	if config.Server.Host == "" {
		return fmt.Errorf("Host is empty")
	}
	if config.Server.Port == "" {
		return fmt.Errorf("PORT is empty")
	}
	return nil
}

