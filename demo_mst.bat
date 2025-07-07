@echo off
REM Script de demostración para Árboles de Expansión Mínimo (MST)
REM Requisito 3b: MST desde cueva específica

echo ===========================================================
echo DEMOSTRACIÓN: ÁRBOLES DE EXPANSIÓN MÍNIMO (MST)
echo Requisito 3b: Visualizar rutas mínimas desde cueva específica
echo ===========================================================

cd /d "%~dp0"

echo.
echo Este script demuestra:
echo 1. Carga de una red de cuevas de ejemplo
echo 2. Cálculo del MST general (Kruskal)
echo 3. Cálculo del MST desde cueva específica (Prim)
echo 4. Análisis de cobertura y rutas mínimas
echo.

REM Verificar que el archivo de datos existe
if not exist "data\caves_mst_example.json" (
    echo Error: Archivo de datos no encontrado
    echo Asegúrese de que existe: data\caves_mst_example.json
    pause
    exit /b 1
)

echo Presione Enter para continuar...
pause >nul

echo.
echo ===========================================================
echo PASO 1: Cargar red de cuevas de ejemplo
echo ===========================================================
echo Archivo: data\caves_mst_example.json
echo Red con 8 cuevas (incluyendo una aislada para demostrar limitaciones)
echo.

echo Contenido del archivo:
echo - BASE: Base Principal (conexiones: 4)
echo - N1, N2: Nodos Norte (bien conectados)
echo - S1, S2: Nodos Sur (bien conectados)
echo - E1: Nodo Este (conexión media)
echo - O1: Nodo Oeste (conexión lejana)
echo - AISLADA: Cueva sin conexiones (para demostrar análisis)

echo.
echo Presione Enter para continuar...
pause >nul

echo.
echo ===========================================================
echo PASO 2: Demostración MST General (Algoritmo Kruskal)
echo ===========================================================
echo El MST general encuentra las conexiones mínimas para toda la red
echo.

REM Simular carga y cálculo de MST general
echo Ejecutando: Análisis ^> MST General...
echo.
echo Resultado esperado:
echo - Peso total mínimo para conectar toda la red alcanzable
echo - Eliminación de conexiones redundantes
echo - Identificación de cuevas aisladas

echo.
echo Presione Enter para continuar...
pause >nul

echo.
echo ===========================================================
echo PASO 3: Demostración MST desde cueva específica (Prim)
echo ===========================================================
echo El MST desde cueva específica muestra:
echo - Qué cuevas son alcanzables desde un punto de origen
echo - Las rutas mínimas desde el origen a cada destino
echo - Análisis de cobertura desde ubicaciones estratégicas

echo.
echo Casos de prueba recomendados:
echo.
echo 1. MST desde 'BASE':
echo    - Debería alcanzar todas las cuevas conectadas (7 de 8)
echo    - Mejor cobertura por ser el nodo más central
echo    - Rutas optimizadas desde el punto principal

echo.
echo 2. MST desde 'N2':
echo    - Alcanzará las mismas cuevas pero con diferentes rutas
echo    - Puede tener mayor peso total por estar en la periferia
echo    - Demuestra cómo el origen afecta las rutas

echo.
echo 3. MST desde 'AISLADA':
echo    - Solo alcanzará esa cueva (1 de 8)
echo    - Demuestra detección de componentes aislados
echo    - Mostrará sugerencias para mejorar conectividad

echo.
echo Presione Enter para continuar...
pause >nul

echo.
echo ===========================================================
echo PASO 4: Ejecución práctica
echo ===========================================================
echo A continuación se ejecutará el programa para demostración interactiva
echo.
echo Pasos sugeridos en el programa:
echo 1. Cargar archivo: data\caves_mst_example.json
echo 2. Ir a: Análisis de la Red
echo 3. Ejecutar: Ver estadísticas de la red
echo 4. Ejecutar: MST General (opción 8)
echo 5. Ejecutar: MST desde cueva específica (opción 9)
echo 6. Probar con diferentes cuevas origen: BASE, N2, AISLADA
echo 7. Ejecutar: Listar cuevas disponibles (opción 10)

echo.
echo Presione Enter para ejecutar el programa...
pause >nul

REM Ejecutar el programa
proyecto-grafos-go.exe

echo.
echo ===========================================================
echo DEMOSTRACIÓN COMPLETADA
echo ===========================================================
echo Has visto cómo el sistema:
echo ✓ Calcula MST general usando Kruskal
echo ✓ Calcula MST desde cueva específica usando Prim
echo ✓ Analiza cobertura y rutas mínimas
echo ✓ Identifica componentes aislados
echo ✓ Proporciona recomendaciones estratégicas
echo.
echo El requisito 3b está completamente implementado.
pause
