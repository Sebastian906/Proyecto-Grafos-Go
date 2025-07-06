#!/bin/bash

# Script para ejecutar las pruebas del proyecto
echo "=== Ejecutando Pruebas del Proyecto ==="

# Navegar al directorio del proyecto
PROJECT_DIR="c:\Users\usuario\Documents\VSCode Proyectos\Proyectos Go\proyecto-grafos-go"
cd "$PROJECT_DIR"

echo "Directorio actual: $(pwd)"
echo ""

echo "=== Prueba de Integración: Requisito 2a ==="
echo "Ejecutando prueba de obstrucción de caminos..."

# Compilar primero para verificar que no hay errores
echo "1. Compilando proyecto..."
go build ./cmd/main.go
if [ $? -eq 0 ]; then
    echo "Compilación exitosa"
else
    echo "Error en compilación"
    exit 1
fi

echo ""
echo "2. Ejecutando prueba de integración..."

# Cambiar a la carpeta de pruebas de integración
cd tests/integration

# Ejecutar la prueba
go test -v .

echo ""
echo "=== Pruebas Completadas ==="
