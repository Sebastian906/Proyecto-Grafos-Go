#!/bin/bash

# Script de prueba para funcionalidad de cuevas inaccesibles
# Proyecto: Sistema de Gestión de Cuevas

echo "=== SCRIPT DE PRUEBA - DETECCION DE CUEVAS INACCESIBLES ==="
echo ""

# Compilar el proyecto
echo "1. Compilando proyecto..."
cd "$(dirname "$0")"
go build ./cmd
if [ $? -ne 0 ]; then
    echo "ERROR: Falló la compilación"
    exit 1
fi
echo "✓ Compilación exitosa"
echo ""

# Ejecutar tests unitarios
echo "2. Ejecutando tests unitarios..."
go test ./tests/unit/accesibilidad_test.go -v
if [ $? -ne 0 ]; then
    echo "ERROR: Fallaron los tests unitarios"
    exit 1
fi
echo "✓ Tests unitarios pasaron"
echo ""

# Ejecutar tests de integración
echo "3. Ejecutando tests de integración..."
go test ./tests/integration/accesibilidad_integration_test.go -v
if [ $? -ne 0 ]; then
    echo "ERROR: Fallaron los tests de integración"
    exit 1
fi
echo "✓ Tests de integración pasaron"
echo ""

echo "=== TODAS LAS PRUEBAS AUTOMATIZADAS PASARON ==="
echo ""
echo "Ahora puedes probar manualmente:"
echo "1. Ejecuta: ./cmd.exe"
echo "2. Ve a 'Gestión de Grafos y Cuevas' (opción 1)"
echo "3. Ve a 'Análisis del grafo' (opción 4)"
echo "4. Prueba 'Detectar cuevas inaccesibles' (opción 4)"
echo "5. Prueba 'Analizar accesibilidad desde cueva específica' (opción 5)"
echo ""
echo "También puedes probar desde:"
echo "- Menú Principal → 'Análisis de Recorridos' (opción 3)"
echo ""

# Mostrar información del grafo cargado
echo "=== INFORMACION DEL GRAFO DE EJEMPLO ==="
echo "Archivo: data/caves_directed_example.json"
echo "Cuevas: 9 cuevas (Silvestre, Tazmania, Coyote, Bunny, Marvin, Piolín, Yayita, Popeye, Correcaminos)"
echo "Conexiones: 13 conexiones dirigidas"
echo "Tipo: Grafo dirigido"
echo ""
echo "Para crear escenarios de prueba interesantes:"
echo "- El grafo ya es dirigido, perfecto para probar cuevas inaccesibles"
echo "- Obstruye conexiones específicas para crear cuevas aisladas"
echo "- Prueba el análisis desde diferentes cuevas de inicio"
echo "- Cambia direcciones de aristas para ver el impacto"
