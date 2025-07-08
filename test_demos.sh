#!/bin/bash

echo "=== VERIFICACIÓN DE TODOS LOS DEMOS ==="
echo ""

demos=(
    "demo_9_cuevas.go"
    "demo_mst.go"
    "demo_prim_simple.go"
    "demo_prim_completo.go"
    "demo_prim_completo_alternativo.go"
    "demo_rutas_acceso.go"
)

total_demos=${#demos[@]}
demos_exitosos=0

for demo in "${demos[@]}"; do
    echo "=========================================="
    echo "Probando: $demo"
    echo "=========================================="
    
    # Verificar compilación
    if go build "demos/$demo" 2>/dev/null; then
        echo "✓ Compilación exitosa"
        
        # Ejecutar demo
        echo "Ejecutando demo..."
        if go run "demos/$demo" > /dev/null 2>&1; then
            echo "✓ Ejecución exitosa"
            demos_exitosos=$((demos_exitosos + 1))
        else
            echo "✗ Error en ejecución"
        fi
    else
        echo "✗ Error de compilación"
    fi
    
    # Limpiar archivos binarios generados
    rm -f demo_9_cuevas demo_mst demo_prim_simple demo_prim_completo demo_prim_completo_alternativo demo_rutas_acceso 2>/dev/null
    
    echo ""
done

echo "=========================================="
echo "RESUMEN DE PRUEBAS"
echo "=========================================="
echo "Total de demos: $total_demos"
echo "Demos exitosos: $demos_exitosos"
echo "Demos fallidos: $((total_demos - demos_exitosos))"

if [ $demos_exitosos -eq $total_demos ]; then
    echo "🎉 ¡TODOS LOS DEMOS FUNCIONAN CORRECTAMENTE!"
    exit 0
else
    echo "⚠️  Algunos demos tienen problemas"
    exit 1
fi
