package config

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

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

func (c *Config) GetDatabaseInformation() string {
	return fmt.Sprintf("name:%s, host:%s, port:%s", c.Database.Database, c.Database.Host, c.Database.Port)
}

func (c *Config) GetCORSConfig() CORSConfig {
	return c.CORS
}

func PrintConfig(cfg *Config, env string) {
	fmt.Println("========== Current Configuration ==========")
	fmt.Printf("Environment      : %s\n", env)

	fmt.Println("[Server]")
	fmt.Printf("Host             : %s\n", cfg.Server.Host)
	fmt.Printf("Port             : %s\n", cfg.Server.Port)
	fmt.Printf("Address          : %s\n", cfg.GetServerAddress())

	fmt.Println("[Static]")
	fmt.Printf("Root Folder      : %s\n", cfg.Static.RootFolder)
	fmt.Printf("Backup Folder    : %s\n", cfg.Static.BackupFolder)

	fmt.Println("[Database]")
	fmt.Printf("Host             : %s\n", emptyFallback(cfg.Database.Host, "(not set)"))
	fmt.Printf("Port             : %s\n", emptyFallback(cfg.Database.Port, "(not set)"))
	fmt.Printf("Username         : %s\n", emptyFallback(cfg.Database.Username, "(not set)"))
	fmt.Printf("Password         : %s\n", maskSecret(cfg.Database.Password))
	fmt.Printf("Database Name    : %s\n", emptyFallback(cfg.Database.Database, "(not set)"))

	fmt.Println("[CORS]")
	fmt.Printf("Allowed Origins  : %s\n", formatSlice(cfg.CORS.AllowedOrigins))
	fmt.Printf("Allowed Methods  : %s\n", formatSlice(cfg.CORS.AllowedMethods))
	fmt.Printf("Allowed Headers  : %s\n", formatSlice(cfg.CORS.AllowedHeaders))

	fmt.Println("===========================================")
}

func validateConfig(config *Config) error {
	var errs []string

	if strings.TrimSpace(config.Server.Host) == "" {
		errs = append(errs, "server.host is empty")
	}
	if err := validatePort("server.port", config.Server.Port); err != nil {
		errs = append(errs, err.Error())
	}

	if strings.TrimSpace(config.Static.RootFolder) == "" {
		errs = append(errs, "static.root_folder is empty")
	}
	if strings.TrimSpace(config.Static.BackupFolder) == "" {
		errs = append(errs, "static.backup_folder is empty")
	}

	if err := validateDatabaseConfig(config.Database); err != nil {
		errs = append(errs, err.Error())
	}

	if err := validateCORSConfig(config.CORS); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "; "))
	}

	return nil
}

func validateDatabaseConfig(db DatabaseConfig) error {
	var errs []string

	hostSet := strings.TrimSpace(db.Host) != ""
	portSet := strings.TrimSpace(db.Port) != ""
	userSet := strings.TrimSpace(db.Username) != ""
	passSet := strings.TrimSpace(db.Password) != ""
	nameSet := strings.TrimSpace(db.Database) != ""

	if !hostSet && !portSet && !userSet && !passSet && !nameSet {
		return nil
	}
	if !hostSet {
		errs = append(errs, "database.host is empty")
	}
	if !portSet {
		errs = append(errs, "database.port is empty")
	} else if err := validatePort("database.port", db.Port); err != nil {
		errs = append(errs, err.Error())
	}
	if !userSet {
		errs = append(errs, "database.username is empty")
	}
	if !nameSet {
		errs = append(errs, "database.database is empty")
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "; "))
	}

	return nil
}

func validateCORSConfig(cors CORSConfig) error {
	var errs []string

	if len(cors.AllowedOrigins) == 0 {
		errs = append(errs, "cors.allowed_origins must not be empty")
	} else {
		for _, origin := range cors.AllowedOrigins {
			origin = strings.TrimSpace(origin)
			if origin == "" {
				errs = append(errs, "cors.allowed_origins contains empty value")
				continue
			}
			if origin != "*" && !isValidOrigin(origin) {
				errs = append(errs, fmt.Sprintf("cors.allowed_origins contains invalid origin: %s", origin))
			}
		}
	}

	if len(cors.AllowedMethods) == 0 {
		errs = append(errs, "cors.allowed_methods must not be empty")
	}

	if len(cors.AllowedHeaders) == 0 {
		errs = append(errs, "cors.allowed_headers must not be empty")
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "; "))
	}

	return nil
}

func validatePort(fieldName, port string) error {
	port = strings.TrimSpace(port)
	if port == "" {
		return fmt.Errorf("%s is empty", fieldName)
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("%s must be a valid number", fieldName)
	}
	if p < 1 || p > 65535 {
		return fmt.Errorf("%s must be between 1 and 65535", fieldName)
	}

	return nil
}

func isValidOrigin(origin string) bool {
	if !strings.Contains(origin, "://") {
		return false
	}

	hostPort := strings.SplitN(origin, "://", 2)
	if len(hostPort) != 2 {
		return false
	}

	scheme := hostPort[0]
	if scheme != "http" && scheme != "https" {
		return false
	}

	host := hostPort[1]
	if host == "" {
		return false
	}

	if strings.Contains(host, "/") {
		return false
	}

	h, p, err := net.SplitHostPort(host)
	if err == nil {
		if strings.TrimSpace(h) == "" {
			return false
		}
		if err := validatePort("origin port", p); err != nil {
			return false
		}
		return true
	}

	return strings.TrimSpace(host) != ""
}

func maskSecret(secret string) string {
	if strings.TrimSpace(secret) == "" {
		return "(not set)"
	}
	return "******"
}

func emptyFallback(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func formatSlice(values []string) string {
	if len(values) == 0 {
		return "(empty)"
	}
	return strings.Join(values, ", ")
}
