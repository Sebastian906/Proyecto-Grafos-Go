#!/bin/bash

# Script para ejecutar pruebas de integración del sistema de grafos
# Proyecto: Sistema de Gestión de Cuevas y Simulación de Camiones

echo "=================================================="
echo "EJECUTOR DE PRUEBAS DE INTEGRACIÓN"
echo "=================================================="
echo

# Función para ejecutar una prueba específica
run_test() {
    local test_name="$1"
    local test_file="$2"
    local description="$3"
    
    echo "EJECUTANDO: Ejecutando: $test_name"
    echo "📝 Descripción: $description"
    echo "📂 Archivo: $test_file"
    echo "INICIANDO: Iniciando en $(date '+%H:%M:%S')"
    echo "--------------------------------------------------"
    
    cd "tests/integration"
    go run "$test_file"
    exit_code=$?
    cd "../.."
    
    echo "--------------------------------------------------"
    if [ $exit_code -eq 0 ]; then
        echo "EXITO: $test_name: EXITOSA"
    else
        echo "ERROR: $test_name: FALLÓ (código: $exit_code)"
    fi
    echo "INICIANDO: Finalizada en $(date '+%H:%M:%S')"
    echo
    
    return $exit_code
}

# Verificar que estamos en el directorio correcto
if [ ! -f "go.mod" ]; then
    echo "ERROR: Error: Este script debe ejecutarse desde la raíz del proyecto"
    echo "   Asegúrese de estar en el directorio que contiene go.mod"
    exit 1
fi

# Verificar que existe la carpeta de datos
if [ ! -d "data" ]; then
    echo "ERROR: Error: No se encontró la carpeta 'data'"
    echo "   Asegúrese de que existe data/caves_example.json"
    exit 1
fi

# Verificar que existe el archivo de datos de ejemplo
if [ ! -f "data/caves_example.json" ]; then
    echo "ERROR: Error: No se encontró data/caves_example.json"
    echo "   Este archivo es necesario para las pruebas de integración"
    exit 1
fi

# Compilar el proyecto principal antes de ejecutar pruebas
echo "COMPILANDO: Compilando proyecto..."
go build ./cmd/main.go
if [ $? -ne 0 ]; then
    echo "ERROR: Error de compilación. Corrija los errores antes de continuar."
    exit 1
fi
echo "EXITO: Compilación exitosa"
echo

# Mostrar información del sistema
echo "INFO: INFORMACIÓN DEL SISTEMA:"
echo "   - Versión de Go: $(go version)"
echo "   - Directorio de trabajo: $(pwd)"
echo "   - Fecha y hora: $(date)"
echo

# Ejecutar pruebas disponibles
total_tests=0
passed_tests=0

# 1. Prueba de simulación de camiones
if [ -f "tests/integration/executable/simulation/main.go" ]; then
    total_tests=$((total_tests + 1))
    echo "EJECUTANDO: Simulación de Camiones"
    echo "DESCRIPCION: Prueba integral del sistema de simulación de camiones con algoritmos DFS y BFS"
    echo "ARCHIVO: executable/simulation/main.go"
    echo "INICIANDO: Iniciando en $(date '+%H:%M:%S')"
    echo "--------------------------------------------------"
    
    cd "tests/integration/executable/simulation"
    go run main.go
    exit_code=$?
    cd "../../../.."
    
    echo "--------------------------------------------------"
    if [ $exit_code -eq 0 ]; then
        echo "EXITO: Simulación de Camiones: EXITOSA"
        passed_tests=$((passed_tests + 1))
    else
        echo "ERROR: Simulación de Camiones: FALLÓ (código: $exit_code)"
    fi
    echo "INICIANDO: Finalizada en $(date '+%H:%M:%S')"
    echo
fi

# 2. Prueba de requerimiento 2a (si existe)
if [ -f "tests/integration/executable/requirement2a/main.go" ]; then
    total_tests=$((total_tests + 1))
    echo "EJECUTANDO: Requerimiento 2a"
    echo "DESCRIPCION: Prueba de funcionalidades de obstrucción de conexiones"
    echo "ARCHIVO: executable/requirement2a/main.go"
    echo "INICIANDO: Iniciando en $(date '+%H:%M:%S')"
    echo "--------------------------------------------------"
    
    cd "tests/integration/executable/requirement2a"
    go run main.go
    exit_code=$?
    cd "../../../.."
    
    echo "--------------------------------------------------"
    if [ $exit_code -eq 0 ]; then
        echo "EXITO: Requerimiento 2a: EXITOSA"
        passed_tests=$((passed_tests + 1))
    else
        echo "ERROR: Requerimiento 2a: FALLÓ (código: $exit_code)"
    fi
    echo "INICIANDO: Finalizada en $(date '+%H:%M:%S')"
    echo
fi

# Resultados finales
echo "=================================================="
echo "RESUMEN DE EJECUCIÓN"
echo "=================================================="
echo "TOTAL: Total de pruebas ejecutadas: $total_tests"
echo "EXITO: Pruebas exitosas: $passed_tests"
echo "ERROR: Pruebas fallidas: $((total_tests - passed_tests))"
echo

if [ $passed_tests -eq $total_tests ]; then
    echo "EXITO: ¡TODAS LAS PRUEBAS PASARON!"
    echo "   El sistema está funcionando correctamente."
    exit 0
else
    echo "ADVERTENCIA:  ALGUNAS PRUEBAS FALLARON"
    echo "   Revise los errores reportados arriba."
    exit 1
fi
