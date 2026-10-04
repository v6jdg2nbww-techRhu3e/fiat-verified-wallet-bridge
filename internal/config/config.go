package config

import (
    "log"
    "os"
)

type Config struct {
    ListenAddr    string
    AdminEmail    string
    AdminPassword string
    Mode          string
}

func Load() *Config {
    cfg := &Config{
        ListenAddr:    getEnv("LISTEN_ADDR", ":8080"),
        AdminEmail:    getEnv("ADMIN_EMAIL", "admin@example.com"),
        AdminPassword: getEnv("ADMIN_PASSWORD", "Admin123!"),
        Mode:          getEnv("APP_MODE", "private"),
    }

    if cfg.AdminEmail == "" || cfg.AdminPassword == "" {
        log.Fatal("admin credentials are required")
    }

    return cfg
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}
