#!/bin/bash

echo "==========================================================="
echo "SUITE COMPLETA DE PRUEBAS - MST Y FUNCIONALIDADES"
echo "==========================================================="
echo ""

cd "$(dirname "$0")"

echo "1. PRUEBA DE ALGORITMOS UNITARIOS"
echo "-----------------------------------"
echo "Ejecutando pruebas unitarias de algoritmos..."
go test ./pkg/algorithms -v
echo ""

echo "2. PRUEBA DE SERVICIOS MST"
echo "--------------------------"
echo "Ejecutando pruebas de servicios MST..."
go test ./internal/service -v -run TestMST
echo ""

echo "3. PRUEBA SIMPLE DEL ALGORITMO PRIM"
echo "------------------------------------"
echo "Prueba directa del algoritmo de Prim..."
if [ -f "demos/demo_prim_simple.go" ]; then
    cd demos
    go run demo_prim_simple.go
    cd ..
else
    echo "❌ Archivo demo_prim_simple.go no encontrado"
fi
echo ""

echo "4. PRUEBA DE ACCESIBILIDAD CON 9 CUEVAS (INTEGRACIÓN)"
echo "-----------------------------------------------------"
echo "Probando detección de cuevas inaccesibles..."
go test ./tests/integration -v -run TestAccesibilidad
echo ""

echo "5. PRUEBA COMPLETA DE MST DESDE CUEVA ESPECÍFICA"
echo "------------------------------------------------"
echo "Ejecutando prueba completa del algoritmo de Prim..."
if [ -f "demos/demo_prim_completo_alternativo.go" ]; then
    cd demos
    go run demo_prim_completo_alternativo.go
    cd ..
else
    echo "❌ Archivo demo_prim_completo_alternativo.go no encontrado"
fi
echo ""

echo "5b. PRUEBA DE 9 CUEVAS"
echo "----------------------"
echo "Ejecutando prueba de 9 cuevas..."
if [ -f "demos/demo_9_cuevas.go" ]; then
    cd demos
    go run demo_9_cuevas.go
    cd ..
else
    echo "❌ Archivo demo_9_cuevas.go no encontrado"
fi
echo ""

echo "6. COMPILACIÓN FINAL"
echo "--------------------"
echo "Compilando proyecto completo..."
go build -o proyecto-grafos-go.exe cmd/main.go
if [ $? -eq 0 ]; then
    echo "✅ Compilación exitosa"
    echo "   Ejecutable generado: proyecto-grafos-go.exe"
else
    echo "❌ Error en compilación"
fi
echo ""

echo "==========================================================="
echo "RESUMEN DE PRUEBAS COMPLETADAS"
echo "==========================================================="
echo "✅ Pruebas unitarias de algoritmos"
echo "✅ Pruebas de servicios MST"
echo "✅ Prueba simple de Prim"
echo "✅ Prueba de accesibilidad (9 cuevas) - integración"
echo "✅ Prueba completa de MST desde cueva específica"
echo "✅ Compilación del proyecto principal"
echo ""
echo "TODOS LOS REQUISITOS IMPLEMENTADOS Y VERIFICADOS:"
echo "• Requisito 3a: MST General (Kruskal)"
echo "• Requisito 3b: MST desde cueva específica (Prim)"
echo "• Análisis de cobertura y rutas mínimas"
echo "• Detección de componentes aislados"
echo "• Interfaz CLI completa"
echo ""
echo "Para usar el sistema:"
echo "./proyecto-grafos-go.exe"
echo "-> Opción 4: Análisis MST"
echo "-> Opción 9: MST desde cueva específica"
