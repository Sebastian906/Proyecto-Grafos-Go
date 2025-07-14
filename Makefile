# Makefile para el Sistema de Gestión de Cuevas

# Variables
APP_NAME := sistema-cuevas
MAIN_FILE := cmd/main.go
BUILD_DIR := build
COVERAGE_DIR := coverage
GO_VERSION := 1.19

# Colores para output
YELLOW := \033[1;33m
GREEN := \033[0;32m
RED := \033[0;31m
NC := \033[0m # No Color

# Detectar sistema operativo
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Windows_NT)
    DETECTED_OS := Windows
    BINARY_EXT := .exe
else ifeq ($(UNAME_S),Linux)
    DETECTED_OS := Linux
    BINARY_EXT := 
else ifeq ($(UNAME_S),Darwin)
    DETECTED_OS := macOS
    BINARY_EXT := 
else
    DETECTED_OS := Unknown
    BINARY_EXT := 
endif

.PHONY: help build clean test coverage run install deps lint format check all

# Mostrar ayuda por defecto
help:
	@echo "$(YELLOW)Sistema de Gestión de Cuevas - Makefile$(NC)"
	@echo ""
	@echo "$(GREEN)Comandos disponibles:$(NC)"
	@echo "  help        - Mostrar esta ayuda"
	@echo "  build       - Compilar la aplicación"
	@echo "  clean       - Limpiar archivos generados"
	@echo "  test        - Ejecutar tests"
	@echo "  coverage    - Generar reporte de cobertura"
	@echo "  run         - Compilar y ejecutar la aplicación"
	@echo "  install     - Instalar dependencias"
	@echo "  deps        - Descargar dependencias"
	@echo "  lint        - Ejecutar linter"
	@echo "  format      - Formatear código"
	@echo "  check       - Verificar código (lint + test)"
	@echo "  all         - Ejecutar todo (deps + lint + test + build)"
	@echo "  cross       - Compilar para múltiples plataformas"
	@echo ""
	@echo "$(GREEN)Información del sistema:$(NC)"
	@echo "  OS detectado: $(DETECTED_OS)"
	@echo "  Go version: $(shell go version 2>/dev/null || echo 'No instalado')"
	@echo ""

# Compilar la aplicación
build:
	@echo "$(YELLOW)Compilando aplicación...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME)$(BINARY_EXT) $(MAIN_FILE)
	@echo "$(GREEN)✓ Aplicación compilada: $(BUILD_DIR)/$(APP_NAME)$(BINARY_EXT)$(NC)"

# Limpiar archivos generados
clean:
	@echo "$(YELLOW)Limpiando archivos generados...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)
	@rm -f $(APP_NAME)$(BINARY_EXT)
	@rm -f *.out
	@rm -f *.exe
	@echo "$(GREEN)✓ Limpieza completada$(NC)"

# Ejecutar tests
test:
	@echo "$(YELLOW)Ejecutando tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)✓ Tests completados$(NC)"

# Generar reporte de cobertura
coverage:
	@echo "$(YELLOW)Generando reporte de cobertura...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | grep total:
	@echo "$(GREEN)✓ Reporte de cobertura generado: $(COVERAGE_DIR)/coverage.html$(NC)"

# Compilar y ejecutar
run: build
	@echo "$(YELLOW)Ejecutando aplicación...$(NC)"
	@./$(BUILD_DIR)/$(APP_NAME)$(BINARY_EXT)

# Instalar dependencias del sistema
install:
	@echo "$(YELLOW)Verificando instalación de Go...$(NC)"
	@which go > /dev/null || (echo "$(RED)✗ Go no está instalado$(NC)" && exit 1)
	@echo "$(GREEN)✓ Go está instalado$(NC)"
	@echo "$(YELLOW)Verificando go mod...$(NC)"
	@test -f go.mod || (echo "$(RED)✗ go.mod no encontrado$(NC)" && exit 1)
	@echo "$(GREEN)✓ go.mod encontrado$(NC)"

# Descargar dependencias
deps:
	@echo "$(YELLOW)Descargando dependencias...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Dependencias actualizadas$(NC)"

# Ejecutar linter (requiere golint)
lint:
	@echo "$(YELLOW)Ejecutando linter...$(NC)"
	@if command -v golint >/dev/null 2>&1; then \
		golint ./...; \
	else \
		echo "$(YELLOW)⚠ golint no está instalado, usando go vet en su lugar$(NC)"; \
		go vet ./...; \
	fi
	@echo "$(GREEN)✓ Linting completado$(NC)"

# Formatear código
format:
	@echo "$(YELLOW)Formateando código...$(NC)"
	@go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "$(YELLOW)⚠ goimports no está instalado$(NC)"; \
	fi
	@echo "$(GREEN)✓ Código formateado$(NC)"

# Verificar código
check: format lint test
	@echo "$(GREEN)✓ Verificación de código completada$(NC)"

# Ejecutar todo el pipeline
all: deps check build
	@echo "$(GREEN)✓ Pipeline completo ejecutado$(NC)"

# Compilar para múltiples plataformas
cross:
	@echo "$(YELLOW)Compilando para múltiples plataformas...$(NC)"
	@mkdir -p $(BUILD_DIR)
	
	@echo "  → Linux AMD64"
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)_linux_amd64 $(MAIN_FILE)
	
	@echo "  → Linux 386"
	@GOOS=linux GOARCH=386 go build -o $(BUILD_DIR)/$(APP_NAME)_linux_386 $(MAIN_FILE)
	
	@echo "  → Windows AMD64"
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)_windows_amd64.exe $(MAIN_FILE)
	
	@echo "  → Windows 386"
	@GOOS=windows GOARCH=386 go build -o $(BUILD_DIR)/$(APP_NAME)_windows_386.exe $(MAIN_FILE)
	
	@echo "  → macOS AMD64"
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)_darwin_amd64 $(MAIN_FILE)
	
	@echo "  → macOS ARM64"
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)_darwin_arm64 $(MAIN_FILE)
	
	@echo "$(GREEN)✓ Compilación multiplataforma completada$(NC)"
	@ls -la $(BUILD_DIR)/

# Ejecutar tests unitarios solamente
test-unit:
	@echo "$(YELLOW)Ejecutando tests unitarios...$(NC)"
	@go test -v ./tests/unit/...
	@echo "$(GREEN)✓ Tests unitarios completados$(NC)"

# Ejecutar tests de integración solamente
test-integration:
	@echo "$(YELLOW)Ejecutando tests de integración...$(NC)"
	@go test -v ./tests/integration/...
	@echo "$(GREEN)✓ Tests de integración completados$(NC)"

# Ejecutar benchmarks
benchmark:
	@echo "$(YELLOW)Ejecutando benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...
	@echo "$(GREEN)✓ Benchmarks completados$(NC)"

# Compilar con optimizaciones de release
release: clean
	@echo "$(YELLOW)Compilando versión de release...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)$(BINARY_EXT) $(MAIN_FILE)
	@echo "$(GREEN)✓ Release compilado: $(BUILD_DIR)/$(APP_NAME)$(BINARY_EXT)$(NC)"

# Verificar dependencias de seguridad
security:
	@echo "$(YELLOW)Verificando dependencias de seguridad...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "$(YELLOW)⚠ gosec no está instalado$(NC)"; \
		echo "  Instalar con: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi
	@echo "$(GREEN)✓ Verificación de seguridad completada$(NC)"

# Crear paquete de distribución
package: release
	@echo "$(YELLOW)Creando paquete de distribución...$(NC)"
	@mkdir -p $(BUILD_DIR)/package
	@cp $(BUILD_DIR)/$(APP_NAME)$(BINARY_EXT) $(BUILD_DIR)/package/
	@cp -r data $(BUILD_DIR)/package/ 2>/dev/null || true
	@cp -r configs $(BUILD_DIR)/package/ 2>/dev/null || true
	@cp docs/README.md $(BUILD_DIR)/package/ 2>/dev/null || true
	@mkdir -p $(BUILD_DIR)/package/logs
	@cd $(BUILD_DIR) && tar -czf $(APP_NAME)_$(DETECTED_OS).tar.gz package/
	@echo "$(GREEN)✓ Paquete creado: $(BUILD_DIR)/$(APP_NAME)_$(DETECTED_OS).tar.gz$(NC)"

# Mostrar información del proyecto
info:
	@echo "$(YELLOW)Información del Proyecto$(NC)"
	@echo "========================"
	@echo "Nombre: Sistema de Gestión de Cuevas"
	@echo "Ejecutable: $(APP_NAME)"
	@echo "Archivo principal: $(MAIN_FILE)"
	@echo "Directorio de build: $(BUILD_DIR)"
	@echo ""
	@echo "$(YELLOW)Estadísticas del código:$(NC)"
	@find . -name "*.go" -not -path "./vendor/*" | wc -l | awk '{print "Archivos Go: " $$1}'
	@find . -name "*.go" -not -path "./vendor/*" -exec cat {} \; | wc -l | awk '{print "Líneas de código: " $$1}'
	@echo ""
	@echo "$(YELLOW)Dependencias:$(NC)"
	@go list -m all 2>/dev/null | head -10

# Configurar hooks de Git
hooks:
	@echo "$(YELLOW)Configurando hooks de Git...$(NC)"
	@mkdir -p .git/hooks
	@echo '#!/bin/sh\nmake check' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "$(GREEN)✓ Hook pre-commit configurado$(NC)"

# Limpiar caché de Go
clean-cache:
	@echo "$(YELLOW)Limpiando caché de Go...$(NC)"
	@go clean -cache
	@go clean -modcache
	@echo "$(GREEN)✓ Caché limpiado$(NC)"

# Mostrar versión
version:
	@echo "$(YELLOW)Versiones del sistema:$(NC)"
	@echo "Go: $(shell go version)"
	@echo "Make: $(shell make --version | head -1)"
	@echo "OS: $(DETECTED_OS)"
