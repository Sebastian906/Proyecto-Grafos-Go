package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config estructura principal de configuración
type Config struct {
	App      AppConfig      `json:"app"`
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Logging  LoggingConfig  `json:"logging"`
}

// AppConfig configuración de la aplicación
type AppConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Debug       bool   `json:"debug"`
}

// DatabaseConfig configuración de base de datos
type DatabaseConfig struct {
	DataDir          string   `json:"data_dir"`
	BackupDir        string   `json:"backup_dir"`
	MaxFileSize      int      `json:"max_file_size_mb"`
	SupportedFormats []string `json:"supported_formats"`
}

// ServerConfig configuración del servidor
type ServerConfig struct {
	Port           int    `json:"port"`
	Host           string `json:"host"`
	ReadTimeout    int    `json:"read_timeout_seconds"`
	WriteTimeout   int    `json:"write_timeout_seconds"`
	MaxConnections int    `json:"max_connections"`
}

// LoggingConfig configuración de logging
type LoggingConfig struct {
	Level      string `json:"level"`
	OutputFile string `json:"output_file"`
	MaxSize    int    `json:"max_size_mb"`
	EnableFile bool   `json:"enable_file"`
}

// DefaultConfig retorna una configuración por defecto
func DefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:        "Sistema de Gestión de Cuevas",
			Version:     "1.0.0",
			Environment: "development",
			Debug:       true,
		},
		Database: DatabaseConfig{
			DataDir:          "data",
			BackupDir:        "backups",
			MaxFileSize:      10,
			SupportedFormats: []string{"json", "xml", "txt"},
		},
		Server: ServerConfig{
			Port:           8080,
			Host:           "localhost",
			ReadTimeout:    30,
			WriteTimeout:   30,
			MaxConnections: 100,
		},
		Logging: LoggingConfig{
			Level:      "info",
			OutputFile: "logs/app.log",
			MaxSize:    50,
			EnableFile: true,
		},
	}
}

// LoadConfig carga la configuración desde un archivo
func LoadConfig(configPath string) (*Config, error) {
	// Si el archivo no existe, crear uno con configuración por defecto
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultConfig()
		if err := SaveConfig(config, configPath); err != nil {
			return nil, fmt.Errorf("error creando archivo de configuración: %w", err)
		}
		return config, nil
	}

	// Leer archivo existente
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo de configuración: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parseando configuración: %w", err)
	}

	// Validar configuración
	if err := ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("configuración inválida: %w", err)
	}

	return &config, nil
}

// SaveConfig guarda la configuración en un archivo
func SaveConfig(config *Config, configPath string) error {
	// Crear directorio si no existe
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("error creando directorio de configuración: %w", err)
	}

	// Serializar configuración
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando configuración: %w", err)
	}

	// Escribir archivo
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo de configuración: %w", err)
	}

	return nil
}

// ValidateConfig valida la configuración
func ValidateConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("configuración no puede ser nil")
	}

	// Validar App
	if config.App.Name == "" {
		return fmt.Errorf("nombre de aplicación no puede estar vacío")
	}
	if config.App.Version == "" {
		return fmt.Errorf("versión de aplicación no puede estar vacía")
	}

	// Validar Database
	if config.Database.DataDir == "" {
		return fmt.Errorf("directorio de datos no puede estar vacío")
	}
	if config.Database.MaxFileSize <= 0 {
		return fmt.Errorf("tamaño máximo de archivo debe ser mayor a 0")
	}

	// Validar Server
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("puerto debe estar entre 1 y 65535")
	}
	if config.Server.Host == "" {
		return fmt.Errorf("host no puede estar vacío")
	}

	// Validar Logging
	validLevels := []string{"debug", "info", "warn", "error"}
	validLevel := false
	for _, level := range validLevels {
		if config.Logging.Level == level {
			validLevel = true
			break
		}
	}
	if !validLevel {
		return fmt.Errorf("nivel de logging debe ser uno de: %v", validLevels)
	}

	return nil
}

// GetDataPath retorna la ruta completa del directorio de datos
func (c *Config) GetDataPath() string {
	return filepath.Clean(c.Database.DataDir)
}

// GetBackupPath retorna la ruta completa del directorio de backups
func (c *Config) GetBackupPath() string {
	return filepath.Clean(c.Database.BackupDir)
}

// GetLogPath retorna la ruta completa del archivo de log
func (c *Config) GetLogPath() string {
	return filepath.Clean(c.Logging.OutputFile)
}

// IsDebug retorna si la aplicación está en modo debug
func (c *Config) IsDebug() bool {
	return c.App.Debug
}

// IsProduction retorna si la aplicación está en modo producción
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// GetServerAddress retorna la dirección completa del servidor
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
