#!/bin/bash

echo "=== PRUEBA RÁPIDA MST DESDE CUEVA ESPECÍFICA ==="
echo "Archivo: caves_mst_example.json"
echo ""

cd "$(dirname "$0")"

# Comprobar archivo
if [ -f "data/caves_mst_example.json" ]; then
    echo "✓ Archivo encontrado: data/caves_mst_example.json"
    echo ""
    echo "Estructura del archivo:"
    echo "- 8 cuevas (7 conectadas + 1 aislada)"
    echo "- BASE: Nodo central con 4 conexiones"
    echo "- N1, N2, S1, S2, E1, O1: Nodos satélite"
    echo "- AISLADA: Sin conexiones (para demostrar análisis)"
    echo ""
    
    echo "Casos de prueba recomendados:"
    echo "1. MST desde 'BASE' -> Mejor cobertura"
    echo "2. MST desde 'N2' -> Nodo periférico"  
    echo "3. MST desde 'AISLADA' -> Componente aislado"
    echo ""
    
    echo "Presione Enter para ejecutar el programa..."
    read
    
    # Ejecutar programa
    ./proyecto-grafos-go.exe
else
    echo "❌ Archivo no encontrado: data/caves_mst_example.json"
    echo "El programa usará los datos por defecto"
    echo ""
    echo "Presione Enter para continuar..."
    read
    ./proyecto-grafos-go.exe
fi
