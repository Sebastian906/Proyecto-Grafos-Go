@echo off
REM Script de prueba para funcionalidad de cuevas inaccesibles
REM Proyecto: Sistema de Gestión de Cuevas

echo === SCRIPT DE PRUEBA - DETECCION DE CUEVAS INACCESIBLES ===
echo.

REM Compilar el proyecto
echo 1. Compilando proyecto...
go build ./cmd
if %errorlevel% neq 0 (
    echo ERROR: Fallo la compilacion
    pause
    exit /b 1
)
echo ✓ Compilacion exitosa
echo.

REM Ejecutar tests unitarios
echo 2. Ejecutando tests unitarios...
go test ./tests/unit/accesibilidad_test.go -v
if %errorlevel% neq 0 (
    echo ERROR: Fallaron los tests unitarios
    pause
    exit /b 1
)
echo ✓ Tests unitarios pasaron
echo.

REM Ejecutar tests de integración
echo 3. Ejecutando tests de integracion...
go test ./tests/integration/accesibilidad_integration_test.go -v
if %errorlevel% neq 0 (
    echo ERROR: Fallaron los tests de integracion
    pause
    exit /b 1
)
echo ✓ Tests de integracion pasaron
echo.

echo === TODAS LAS PRUEBAS AUTOMATIZADAS PASARON ===
echo.
echo Ahora puedes probar manualmente:
echo 1. Ejecuta: cmd.exe
echo 2. Ve a 'Gestion de Grafos y Cuevas' (opcion 1^)
echo 3. Ve a 'Analisis del grafo' (opcion 4^)
echo 4. Prueba 'Detectar cuevas inaccesibles' (opcion 4^)
echo 5. Prueba 'Analizar accesibilidad desde cueva especifica' (opcion 5^)
echo.
echo Tambien puedes probar desde:
echo - Menu Principal → 'Analisis de Recorridos' (opcion 3^)
echo.

REM Mostrar información del grafo cargado
echo === INFORMACION DEL GRAFO DE EJEMPLO ===
echo Archivo: data/caves_directed_example.json
echo Cuevas: 9 cuevas (Silvestre, Tazmania, Coyote, Bunny, Marvin, Piolin, Yayita, Popeye, Correcaminos^)
echo Conexiones: 13 conexiones dirigidas
echo Tipo: Grafo dirigido
echo.
echo Para crear escenarios de prueba interesantes:
echo - El grafo ya es dirigido, perfecto para probar cuevas inaccesibles
echo - Obstruye conexiones especificas para crear cuevas aisladas
echo - Prueba el analisis desde diferentes cuevas de inicio
echo - Cambia direcciones de aristas para ver el impacto
echo.

pause
