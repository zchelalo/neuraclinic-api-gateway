package bootstrap

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var (
	config   Config
	configMu sync.RWMutex
)

type GRPCConfig struct {
	Addr               string
	TLSEnabled         bool
	CACertPath         string
	InsecureSkipVerify bool
}

type Config struct {
	Environment    string
	ServiceName    string
	Port           int
	AllowedOrigins []string

	AuthGRPC           GRPCConfig
	UsersGRPC          GRPCConfig
	RecordsGRPC        GRPCConfig
	LocationGRPC       GRPCConfig
	FileManagementGRPC GRPCConfig

	InternalServiceToken string

	CookieDomain      string
	CookieSecure      bool
	AccessCookieName  string
	RefreshCookieName string
}

func LoadConfig(dotenvPath string) (Config, error) {
	if dotenvPath != "" {
		_ = godotenv.Load(dotenvPath)
	}

	cfg := Config{
		Environment:    getEnv("ENVIRONMENT", "development"),
		ServiceName:    getEnv("SERVICE_NAME", "neuraclinic-api-gateway"),
		Port:           getEnvInt("PORT", 8000),
		AllowedOrigins: splitCSV(getEnv("ALLOWED_ORIGINS", "")),
		AuthGRPC: GRPCConfig{
			Addr:               getEnv("AUTH_GRPC_ADDR", ""),
			TLSEnabled:         getEnvBool("AUTH_GRPC_TLS_ENABLED", true),
			CACertPath:         getEnv("AUTH_GRPC_CA_CERT_PATH", ""),
			InsecureSkipVerify: getEnvBool("AUTH_GRPC_INSECURE_SKIP_VERIFY", false),
		},
		UsersGRPC: GRPCConfig{
			Addr:               getEnv("USERS_GRPC_ADDR", ""),
			TLSEnabled:         getEnvBool("USERS_GRPC_TLS_ENABLED", true),
			CACertPath:         getEnv("USERS_GRPC_CA_CERT_PATH", ""),
			InsecureSkipVerify: getEnvBool("USERS_GRPC_INSECURE_SKIP_VERIFY", false),
		},
		RecordsGRPC: GRPCConfig{
			Addr:               getEnv("RECORDS_GRPC_ADDR", ""),
			TLSEnabled:         getEnvBool("RECORDS_GRPC_TLS_ENABLED", true),
			CACertPath:         getEnv("RECORDS_GRPC_CA_CERT_PATH", ""),
			InsecureSkipVerify: getEnvBool("RECORDS_GRPC_INSECURE_SKIP_VERIFY", false),
		},
		LocationGRPC: GRPCConfig{
			Addr:               getEnv("LOCATION_GRPC_ADDR", ""),
			TLSEnabled:         getEnvBool("LOCATION_GRPC_TLS_ENABLED", true),
			CACertPath:         getEnv("LOCATION_GRPC_CA_CERT_PATH", ""),
			InsecureSkipVerify: getEnvBool("LOCATION_GRPC_INSECURE_SKIP_VERIFY", false),
		},
		FileManagementGRPC: GRPCConfig{
			Addr:               getEnv("FILE_MANAGEMENT_GRPC_ADDR", ""),
			TLSEnabled:         getEnvBool("FILE_MANAGEMENT_GRPC_TLS_ENABLED", true),
			CACertPath:         getEnv("FILE_MANAGEMENT_GRPC_CA_CERT_PATH", ""),
			InsecureSkipVerify: getEnvBool("FILE_MANAGEMENT_GRPC_INSECURE_SKIP_VERIFY", false),
		},
		InternalServiceToken: getEnv("INTERNAL_SERVICE_TOKEN", ""),
		CookieDomain:         getEnv("COOKIE_DOMAIN", ""),
		CookieSecure:         getEnvBool("COOKIE_SECURE", false),
		AccessCookieName:     getEnv("ACCESS_COOKIE_NAME", "access_token"),
		RefreshCookieName:    getEnv("REFRESH_COOKIE_NAME", "refresh_token"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	setConfig(cfg)
	return cfg, nil
}

func GetConfig() Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return config
}

func setConfig(cfg Config) {
	configMu.Lock()
	config = cfg
	configMu.Unlock()
}

func (c Config) Validate() error {
	required := map[string]string{
		"AUTH_GRPC_ADDR":            c.AuthGRPC.Addr,
		"USERS_GRPC_ADDR":           c.UsersGRPC.Addr,
		"RECORDS_GRPC_ADDR":         c.RecordsGRPC.Addr,
		"LOCATION_GRPC_ADDR":        c.LocationGRPC.Addr,
		"FILE_MANAGEMENT_GRPC_ADDR": c.FileManagementGRPC.Addr,
		"INTERNAL_SERVICE_TOKEN":    c.InternalServiceToken,
		"ACCESS_COOKIE_NAME":        c.AccessCookieName,
		"REFRESH_COOKIE_NAME":       c.RefreshCookieName,
	}

	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("missing required config key: %s", key)
		}
	}
	if c.Port <= 0 {
		return fmt.Errorf("PORT must be greater than zero")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
