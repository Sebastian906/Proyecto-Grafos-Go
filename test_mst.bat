@echo off
echo ================================================
echo PRUEBA RÁPIDA: MST desde Cueva Específica
echo ================================================
echo.
echo Este script demuestra cómo probar la funcionalidad:
echo.
echo 1. Cargar archivo: data/caves_mst_example.json
echo 2. Ir a: Análisis de la Red
echo 3. Probar MST desde diferentes cuevas
echo.
echo ================================================
echo CASOS DE PRUEBA RECOMENDADOS:
echo ================================================
echo.
echo Caso 1: MST desde 'BASE'
echo - Debería alcanzar 7 de 8 cuevas
echo - Peso total alrededor de 41.0
echo - Cobertura: 87.5%%
echo.
echo Caso 2: MST desde 'N2' (nodo periférico)
echo - Debería alcanzar 7 de 8 cuevas
echo - Peso total mayor que desde BASE
echo - Demuestra efecto del punto de origen
echo.
echo Caso 3: MST desde 'AISLADA'
echo - Solo alcanzará 1 cueva (sí misma)
echo - Cobertura: 12.5%%
echo - Demuestra detección de componentes aislados
echo.
echo ================================================
echo Presione cualquier tecla para ejecutar...
pause > nul

proyecto-grafos-go.exe

echo.
echo ================================================
echo PRUEBA COMPLETADA
echo ================================================
echo.
echo ¿Funcionó correctamente? Deberías haber visto:
echo ✓ Carga exitosa del archivo JSON
echo ✓ Menú de análisis con opciones 9 y 10
echo ✓ Cálculo de MST desde diferentes cuevas
echo ✓ Análisis de cobertura y rutas mínimas
echo.
pause
