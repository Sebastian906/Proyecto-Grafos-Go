@echo off
REM Script para ejecutar pruebas de integración del siif existREM Ejecutar prueba de simulación de camiones
    set exit_code=!errorlevel!
    cd ..\..\..\..
    
    echo --------------------------------------------------
    if !exit_code! equ 0 (
        echo EXITO: Requerimiento 2a: EXITOSA
        set /a passed_tests+=1
    ) else (
        echo ERROR: Requerimiento 2a: FALLÓ ^(código: !exit_code!^)
    )tests\integration\executable\simulation\main.go" (
    set /a total_tests+=1
    echo EJECUTANDO: Simulación de Camiones
    echo DESCRIPCION: Prueba integral del sistema de simulación de camiones con algoritmos DFS y BFS
    echo ARCHIVO: executable\simulation\main.go
    echo INICIANDO: Iniciando en %TIME%
    echo --------------------------------------------------
    
    cd "tests\integration\executable\simulation"
    go run main.goegration\executable\simulation_integration.go" (
    set /a total_tests+=1
    echo EJECUTANDO: Simulación de Camiones
    echo DESCRIPCION: Prueba integral del sistema de simulación de camiones con algoritmos DFS y BFS
    echo ARCHIVO: executable\simulation_integration.go
    echo INICIANDO: Iniciando en %time%
    echo --------------------------------------------------
    
    cd "tests\integration\executable"
    go run simulation_integration.gorafos
if exist "tests\integration\executable\requirement2a\main.go" (
    set /a total_tests+=1
    echo EJECUTANDO: Requerimiento 2a
    echo DESCRIPCION: Prueba de funcionalidades de obstrucción de conexiones
    echo ARCHIVO: executable\requirement2a\main.go
    echo INICIANDO: Iniciando en %time%
    echo --------------------------------------------------
    
    cd "tests\integration\executable\requirement2a"
    go run main.goto: Sistema de Gestión de Cuevas y Simulación de Camiones

echo ==================================================
echo EJECUTOR DE PRUEBAS DE INTEGRACIÓN
echo ==================================================
echo.

REM Verificar que estamos en el directorio correcto
if not exist "go.mod" (
    echo ERROR:: Este script debe ejecutarse desde la raíz del proyecto
    echo    Asegúrese de estar en el directorio que contiene go.mod
    pause
    exit /b 1
)

REM Verificar que existe la carpeta de datos
if not exist "data" (
    echo ERROR:: No se encontró la carpeta 'data'
    echo    Asegúrese de que existe data/caves_example.json
    pause
    exit /b 1
)

REM Verificar que existe el archivo de datos de ejemplo
if not exist "data\caves_example.json" (
    echo ERROR:: No se encontró data\caves_example.json
    echo    Este archivo es necesario para las pruebas de integración
    pause
    exit /b 1
)

REM Compilar el proyecto antes de ejecutar pruebas
echo 🔨 Compilando proyecto...
go build ./...
if errorlevel 1 (
    echo ERROR: de compilación. Corrija los errores antes de continuar.
    pause
    exit /b 1
)
echo EXITO: Compilación exitosa
echo.

REM Mostrar información del sistema
echo 📊 INFORMACIÓN DEL SISTEMA:
go version
echo    - Directorio de trabajo: %CD%
echo    - Fecha y hora: %DATE% %TIME%
echo.

set total_tests=0
set passed_tests=0

REM Ejecutar prueba de simulación de camiones
if exist "tests\integration\simulation_integration.go" (
    set /a total_tests+=1
    echo 🧪 Ejecutando: Simulación de Camiones
    echo 📝 Descripción: Prueba integral del sistema de simulación de camiones con algoritmos DFS y BFS
    echo 📂 Archivo: simulation_integration.go
    echo INICIANDO: Iniciando en %TIME%
    echo --------------------------------------------------
    
    cd tests\integration
    go run simulation_integration.go
    set exit_code=!errorlevel!
    cd ..\..\..\..
    
    echo --------------------------------------------------
    if !exit_code! equ 0 (
        echo EXITO: Simulación de Camiones: EXITOSA
        set /a passed_tests+=1
    ) else (
        echo ERROR: Simulación de Camiones: FALLÓ ^(código: !exit_code!^)
    )
    echo INICIANDO: Finalizada en %TIME%
    echo.
)

REM Ejecutar prueba de requerimiento 2a (si existe)
if exist "tests\integration\requirement_2a_test.go" (
    set /a total_tests+=1
    echo 🧪 Ejecutando: Requerimiento 2a
    echo 📝 Descripción: Prueba de funcionalidades de obstrucción de conexiones
    echo 📂 Archivo: requirement_2a_test.go
    echo INICIANDO: Iniciando en %TIME%
    echo --------------------------------------------------
    
    cd tests\integration
    go run requirement_2a_test.go
    set exit_code=!errorlevel!
    cd ..\..
    
    echo --------------------------------------------------
    if !exit_code! equ 0 (
        echo EXITO: Requerimiento 2a: EXITOSA
        set /a passed_tests+=1
    ) else (
        echo ❌ Requerimiento 2a: FALLÓ ^(código: !exit_code!^)
    )
    echo INICIANDO: Finalizada en %TIME%
    echo.
)

REM Resultados finales
echo ==================================================
echo 📋 RESUMEN DE EJECUCIÓN
echo ==================================================
echo 🧪 Total de pruebas ejecutadas: %total_tests%
echo EXITO: Pruebas exitosas: %passed_tests%
set /a failed_tests=%total_tests%-%passed_tests%
echo ❌ Pruebas fallidas: %failed_tests%
echo.

if %passed_tests% equ %total_tests% (
    echo 🎉 ¡TODAS LAS PRUEBAS PASARON!
    echo    El sistema está funcionando correctamente.
    echo.
    echo Presione cualquier tecla para continuar...
    pause >nul
    exit /b 0
) else (
    echo ADVERTENCIA:  ALGUNAS PRUEBAS FALLARON
    echo    Revise los errores reportados arriba.
    echo.
    echo Presione cualquier tecla para continuar...
    pause >nul
    exit /b 1
)
