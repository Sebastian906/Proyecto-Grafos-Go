@echo off
echo ===================================================
echo TESTS DE ALGORITMO PRIM - MST DESDE CUEVA ESPECIFICA
echo ===================================================
echo.

cd /d "%~dp0"

echo Ejecutando prueba completa del algoritmo de Prim...
echo.
cd tests\manual
go run test_prim_completo.go

echo.
echo ===================================================
echo Ejecutando pruebas unitarias...
echo ===================================================
cd ..\..
go test ./pkg/algorithms -v -run=".*Prim.*"

echo.
echo ===================================================
echo TODAS LAS PRUEBAS COMPLETADAS
echo ===================================================
echo.
echo Los archivos de prueba estan en:
echo - tests\manual\test_prim_completo.go (prueba completa)
echo - tests\manual\test_prim_directo.go (referencia)
echo - pkg\algorithms\prim_test.go (pruebas unitarias)
echo.
pause
