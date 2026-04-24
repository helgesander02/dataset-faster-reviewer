package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfigForEnvironment(env string) (*Config, error) {
	if err := initializeViper(env); err != nil {
		return nil, err
	}

	config, err := unmarshalConfig()
	if err != nil {
		return nil, err
	}

	overrideWithEnvironmentVariables(config)

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("configuration verification failed: %w", err)
	}

	return config, nil
}

func initializeViper(env string) error {
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
			fmt.Printf("Configuration file %s.yaml not found, using defaults\n", configName)
			return nil
		}
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	return nil
}

func unmarshalConfig() (*Config, error) {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}
	return &config, nil
}

func overrideWithEnvironmentVariables(config *Config) {
	if allowedOrigins := os.Getenv("ALLOWED_ORIGINS"); allowedOrigins != "" {
		origins := strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
		config.CORS.AllowedOrigins = origins
	}
}

func setDefaults() {
	setServerDefaults()
	setStaticDefaults()
	setDatabaseDefaults()
	setCORSDefaults()
}

func setServerDefaults() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "localhost")
}

func setStaticDefaults() {
	viper.SetDefault("static.root_folder", "./static")
	viper.SetDefault("static.backup_folder", "./static")
}

func setDatabaseDefaults() {
	viper.SetDefault("database.host", "")
	viper.SetDefault("database.port", "")
	viper.SetDefault("database.username", "")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.database", "")
}

func setCORSDefaults() {
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3000"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"*"})
}
