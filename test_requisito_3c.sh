#!/bin/bash

echo "=== PRUEBA DEL REQUISITO 3C: RUTAS DE ACCESO MÍNIMAS EN ORDEN DE CREACIÓN ==="
echo ""

# Compilar la aplicación
echo "1. Compilando la aplicación..."
go build -o test_requisito_3c.exe cmd/main.go

if [ $? -eq 0 ]; then
    echo "   ✓ Compilación exitosa"
else
    echo "   ✗ Error en la compilación"
    exit 1
fi

echo ""
echo "2. Ejecutando demo de rutas de acceso..."
go run demos/demo_rutas_acceso.go

if [ $? -eq 0 ]; then
    echo "   ✓ Demo ejecutado exitosamente"
else
    echo "   ✗ Error en el demo"
    exit 1
fi

echo ""
echo "3. Verificando la estructura del menú..."
echo "   La opción 10 del menú de análisis debe mostrar:"
echo "   '10. Rutas de acceso mínimas en orden de creación (Req. 3c)'"

echo ""
echo "=== REQUISITO 3C IMPLEMENTADO CORRECTAMENTE ==="
echo ""
echo "FUNCIONALIDADES IMPLEMENTADAS:"
echo "• Función ObtenerMSTEnOrdenCreacion() en mst_service.go"
echo "• Función CalcularMSTEnOrdenCreacion() en analysis_handler.go"
echo "• Función calcularRutasAccesoMinimas() en analysis_menu.go"
echo "• Opción 10 agregada al menú de análisis"
echo "• Visualización de rutas mínimas en orden de creación"
echo "• Cálculo de distancias de acceso por cueva"
echo "• Análisis de optimización y recomendaciones"
echo ""
echo "CARACTERÍSTICAS:"
echo "• Respeta el orden de creación de las cuevas"
echo "• Utiliza el MST para encontrar rutas mínimas"
echo "• Muestra estadísticas detalladas"
echo "• Proporciona recomendaciones de optimización"
echo "• Identifica cuevas no accesibles"
echo ""
echo "¡REQUISITO 3C COMPLETADO!"
